package content

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	CreateCmdUse   string = "create [subcommand]"
	CreateCmdShort string = "Cria o cluster Girus"
)

const (
	CreateClusterCmdUse   string = "cluster"
	CreateClusterCmdShort string = "Cria um cluster Girus"
	CreateClusterCmdLong  string = "Cria um cluster Kind com o nome \"girus\" e implanta todos os componentes necessários. Por padrão, o deployment embutido no binário é utilizado."
)

const (
	DockerMacOSinstructions = `
Para macOS, recomendamos usar Colima (alternativa leve ao Docker Desktop):
1. Instale o Homebrew caso não tenha:
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
2. Instale o Colima e o Docker CLI:
  brew install colima docker
3. Inicie o Colima:
  colima start
Alternativamente, você pode instalar o Docker Desktop para macOS de:
https://www.docker.com/products/docker-desktop
`

	DockerLinuxInstructions = `
1. Para Linux, use o script de instalação oficial:
  curl -fsSL https://get.docker.com | bash
2. Após a instalação, adicione seu usuário ao grupo docker para evitar usar sudo:
  sudo usermod -aG docker $USER
  newgrp docker
3. Inicie o serviço:
  sudo systemctl enable docker
  sudo systemctl start docker
`
	PodmanMacOSInstructions = `
Para macOS, recomendamos Podman:
1. Instale o Homebrew caso não tenha:
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
2. Instale o Podman
  brew install podman
3. Inicie o Podman:
  podman machine init
  podman machine start
`

	PodmanLinuxInstructions = `
1. Para Linux, use o script de instalação oficial:
  curl -fsSL https://get.docker.com | bash
2. E inicie o serviço:
  sudo systemctl enable podman
  sudo systemctl start podman
3. Opicional: Após a instalação, para utilizar podman, rootless evitando sudo:
  Siga as instruções do site oficial:
  https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md
`

	PodmanWindowsInstructions = `
Visite https://github.com/containers/podman/blob/main/docs/tutorials/podman-for-windows.md para instruções de instalação para seu sistema operacional
`

	OtherOSInstructions = `
Visite https://www.docker.com/products/docker-desktop para instruções de instalação para seu sistema operacional
`
)

const (
	DockerInfoDarwin = `
Para macOS com Colima:
  colima start
Para Docker Desktop:
  Inicie o aplicativo Docker Desktop
`

	PodmanInfoDarwin = `
Para Podman:
	Inicie a machine com: podman machine start
`

	DockerInfoLinux = `
Inicie o serviço Docker:
  sudo systemctl start docker
`

	PodmanInfoLinux = `
Para Podman:
  sudo systemctl start podman
`
)

func CreateHeader() string {
	headerColor := color.New(color.FgCyan, color.Bold).SprintFunc()

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		strings.Repeat("─", 80),
		headerColor("GIRUS CREATE"),
		strings.Repeat("─", 80),
		headerColor("Verificando atualizações..."),
	)
}
