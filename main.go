package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bitfield/script"
)

func main() {
	// https://go.dev/doc/install

	update := flag.Bool("system", false, "update system")
	docker := flag.Bool("docker", false, "install docker engine")
	gitlabRunner := flag.Bool("gitlab", false, "install gitlab runner")

	flag.Parse()

	if *update {
		updateSystem()
	} else if *docker {
		installDocker()
	} else if *gitlabRunner {
		installGitlabRunner()
	} else {
		fmt.Println("run the command with --help")
	}
}

func updateSystem() {
	// update kernel
	_, err := script.Exec("sudo apt-get update").Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("sudo apt-get upgrade --yes").Stdout()
	if err != nil {
		return
	}
}

func installDocker() {
	// install docker engine
	_, err := script.Exec("sudo apt-get install --yes apt-transport-https ca-certificates curl gnupg lsb-release").Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("curl -fsSL https://download.docker.com/linux/ubuntu/gpg").Exec("sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg --yes").Stdout()
	if err != nil {
		return
	}
	releaseModel, err := script.Exec("lsb_release -cs").String()
	responseApp := fmt.Sprintf("echo deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu %s stable", releaseModel)
	_, err = script.Exec(responseApp).Exec("sudo tee /etc/apt/sources.list.d/docker.list > /dev/null").Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("sudo apt-get update").Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("sudo apt-get install --yes docker-ce docker-ce-cli containerd.io docker-compose-plugin").Stdout()
	if err != nil {
		return
	}
	user := os.Getenv("USER")
	usermod := fmt.Sprintf("sudo usermod -aG docker %s", user)
	_, err = script.Exec(usermod).Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("docker run hello-world").Stdout()
	if err != nil {
		return
	}
	//

	script.Echo("\nThe process has been successfully completed.\n").Stdout()
}

func installGitlabRunner() {
	// install gitlab-runner
	_, err := script.Exec("gitlab-runner stop").Stdout()
	_, err = script.Exec("sudo curl -L --output /usr/local/bin/gitlab-runner 'https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64'").Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("sudo chmod +x /usr/local/bin/gitlab-runner").Stdout()
	if err != nil {
		return
	}
	s, err := script.Exec("sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash").String()
	if err != nil {
		if s != "useradd: user 'gitlab-runner' already exists\n" {
			return
		}
	}
	s, err = script.Exec("sudo gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner").String()
	if err != nil {
		if s1, _ := script.Echo(s).Match("gitlab-runner: Init already exists").String(); len(s1) <= 0 {
			return
		}
	}
	_, err = script.Exec("sudo usermod -aG docker gitlab-runner").Stdout()
	if err != nil {
		return
	}
	_, err = script.Exec("systemctl restart docker").Stdout()
	if err != nil {
		return
	}
	_, err = script.Echo(" ").WriteFile("/home/gitlab-runner/.bash_logout")
	if err != nil {
		return
	}
	_, err = script.Exec("gitlab-runner start").Stdout()
	if err != nil {
		return
	}
	//

	// finish
	script.Echo("\nThe process has been successfully completed.\n").Stdout()
}
