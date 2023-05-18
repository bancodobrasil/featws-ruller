[![Go Reference](https://pkg.go.dev/badge/github.com/abu-lang/goabu.svg)](https://pkg.go.dev/github.com/bancodobrasil/featws-ruller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/bancodobrasil/featws-ruller/blob/develop/LICENSE)


# **Featws Ruller** [![About_de](https://github.com/yammadev/flag-icons/blob/master/png/BR.png?raw=true)](https://github.com/bancodobrasil/featws-transpiler/blob/develop/README_ptbr.md)

- O projeto featws-ruller é uma implementação do [grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine), usado para aveliar planilhas de regaras (.grl).

## Software Necessário
- Será necessário ter instalado em sua máquina a **Linguagem de Programação Go** para rodar o projeto. Você pode fazer o download na pádina oficial [aqui](https://go.dev/doc/install).
- Clone o repositório **featws-transpiler** para a sua máquina local e tenha verifique se o projeto transpiler e ruller estão na mesma pasta.Você pode achar o repositorio do featws-transpiler [aqui](https://github.com/bancodobrasil/featws-transpiler).

## Inicializando o Projeto
- Clone o projeto para sua máquina local.
- Com a pasta do projeto aberto (*../featws-ruller/main.go*), abra o arquivo  _main.go_ e o terminal integrado, digite o comando `go run main.go`. Se voce utiliza o sistema OS ou windows, voce tambem pode dar build e executar o projeto com os comandos: `go build && ./featws-ruler.exe`, caso use windows, ou  `go build -o ruller && ./ruller $@` se utiliza Mac.

## Testando diferentes folhas de regras
- Verifique se voce possui em seu workspace o **featws-transpiler** e copie o caminho do arquivo .grl para o novo caso a ser testado. Você pode encontrar isso nos casos _tests_ -> cases.
- Agora basta substituir a variável env "FEATWS_RULLER_DEFAULT_RULES" no arquivo .env na regra, pelo novo caminho, e executar conforme as instruções acima.

## Folha de regras de teste com resolvers
- Para testar se o resolver está carregado, você deve definir a URL **featws-resolver-bridge** no arquivo .env.

## Carregando uma folha de regras de uma fonte remota
- Para carregar uma planilha de uma fonte remota, basta alterar a variável .env "FEATWS_RULLER_RESOURCE_LOADER_URL" apontada para sua URL.

# Usando principais endpoints 
_Por padrão a porta utilizada será a :8000_
- GET **http://localhost:SUAPORTAESCOLHIDA/**
  - Retornará uma mensagem simples ao cliente, como: "FeatWS Ruller Works!!!"

- POST **http://localhost:SUAPORTAESCOLHIDA/api/v1/eval**
  - Neste ponto final você deve ter que passar um corpo, que são os parâmetros definidos pela pasta rulesheet no featws-transpiler. Usando o case 0001, por exemplo, o corpo deve ser:
    ```json
        {
            "mynumber": "45"
        }
    ```
   - Após a solicitação ter sido enviada, a resposta deve ser algo assim: (essa resposta é definida pelo arquivo .featws na pasta ruleshet, nesse caso é false porque a condição é meunúmero> 12)
        ```json 
            {
                "myboolfeat": false
            }
        ```
- GET **http://localhost:SUAPORTAESCOLHIDA/swagger/index.html**
    - No seu navegador, você pode ver a documentação do swagger da api.

- GET **http://localhost:SUAPORTAESCOLHIDA/health/live?full=1** 
    - Este endpoint verificará o status ativo do aplicação:

- GET **http://localhost:SUAPORTAESCOLHIDA/health/ready?full=1**
    - Este endpoint verificará o status se o serviços está pronto para ser usado ​​pelo projeto da ruller.
    
    ```json
    {
    "goroutine-threshold": "OK",
    "resolver-bridge": "OK",
    "resource-loader": "OK"
    }
    ```