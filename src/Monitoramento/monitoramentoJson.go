package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const monitoramento = 3
const delay = 1

type SitesInterface struct {
	Sites []string `json: "sites"`
}

func main() {
	exibirIntroducao()
	for {
		exibeMenu()

		comando := lerComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa.")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func exibirIntroducao() {
	nome := "Doug"
	versao := 1.1
	fmt.Println("Hello Mr.", nome)
	fmt.Println("this program is in version", versao)
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do Programa")
}

func lerComando() int {
	var comando int
	fmt.Scan(&comando)
	fmt.Println("O comando escolhido foi", comando)

	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := lerSitesDoJson()

	for count := 0; count < monitoramento; count++ {
		for posicaoSite, site := range sites {
			fmt.Println("Testando site", posicaoSite+1, ":", site)
			testarSite(site)
			//Delay de resposta
			time.Sleep(delay * time.Second)
		}
	}
}

func testarSite(site string) {
	resposta, err := http.Get(site)

	//Verifica se ocorreu algum erro ao conectar com site
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resposta.StatusCode == 200 {
		fmt.Println("Site:", site, "Foi carregado com sucesso")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status Code:", resposta.StatusCode)
		registraLog(site, false)
	}
}

func lerSitesDoJson() []string {
	var sites []string
	jsonFile, err := os.Open("sites.json")

	//Verifica se ocorreu algum erro ao conectar com json
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	byteValueJSON, _ := ioutil.ReadAll(jsonFile)
	objSites := SitesInterface{}
	err = json.Unmarshal(byteValueJSON, &objSites)

	if err != nil {
		fmt.Println(err)
	}

	//Lê o json, e caso exista sites, ele insere no array
	for _, site := range objSites.Sites {

		sites = append(sites, site)

		//Verifica se esta no final da linha
		if err == io.EOF {
			break
		}
	}

	defer jsonFile.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + "- online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(delay * time.Second)
	fmt.Println(string(arquivo))
}
