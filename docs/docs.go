// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/eval": {
            "post": {
                "security": [
                    {
                        "Authentication Api Key": []
                    }
                ],
                "description": "Para a realiza os testes basta clicar em *Try it out*, complete a folha de regra com os dados desejados e os demais campos caso necessário, em seguida, clique em *Execute*.\n\nA seguir, serão apresentados alguns exemplos de testes:\n\n[**Exemplo 1**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0001%20-%20one_parameter) **Testando uma variável**:\n` + "`" + `` + "`" + `` + "`" + `\n\"myboolfeat\" = \"$mynumber \u003c 12\"\n` + "`" + `` + "`" + `` + "`" + `\n\nNesse exemplo, é possível testar a feature *myboolfeat*. Ao abrir o arquivo [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0001%20-%20one_parameter/rules.featws), é possível observar que, se o valor de *mynumber* for menor que 12, a feature *myboolfeat* retornará *true*. Caso contrário, se for maior ou igual a 12, o retorno será *false*. Portanto, para testar essa regra, basta inserir o seguinte corpo de *Parameters*.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"mynumber\": \"1\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 2**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0003%20-%20simple_group) **Regra com um grupo**:\n\nNesse exemplo vamos testar a utilização de um grupo espeífico. Ao enviarmos um **clientingroup** que tenha no grupo [mygroup](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0003%20-%20simple_group/groups/mygroup.json), espera-se como retorno **mygroup** = true, caso seja passado uma agencia e conta válidas como a seguir:\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"branch\": \"00000\",\n\"account\": \"00000000\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 3**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0004%20-%20default_value_true) **Regra com um valor padrão**:\n\nNesse exemplo vamos testar a utilização de um valor padrão. Ao enviarmos um **gender** F sua resposta esperada será \"female\" = \"true\", e \"male\" =\"false\", pois no [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/rules.featws) da regra diz que female é o inverso de male. Logo caso seja enviado o valor de \"gender\": \"M\", teremos female =\"false\" e male =\"true\". Sendo enviado qualque outra letra como parâmetro será recebido female=\"false\" e male=\"true\", pois o default male está declarado primeiro do que o female no arquivo [features.json](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/features.json) da regra.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"gender\": \"F\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 4**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0005%20-%20precedence_order) **Regra com ordem de procedência**:\n\nNesse exemplo vamos testar a utilização de uma ordem de procedência. Similar ao caso anterior, se é enviado um valor menor ou igual à 18, receberemos como resposta \"menor_de_idade\": \"true\" e \"maior_de_idade\": \"false\", caso seja enviado um valor maior que 18 deveremos receber \"menor_de_idade\": \"false\" e \"maior_de_idade\": \"true\". Caso não seja enviado nenhum valor de idade, será interpretado como \"menor_de_idade\": \"true\" e \"maior_de_idade\":\"false\".\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"idade\": \"21\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 5**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0008%20-%20group%20_intersection) **Regra com interseção de grupos**:\n\nA interseção de grupos ocorre quando dois ou mais grupos têm elementos em comum, ou seja, há um conjunto de elementos que pertencem a todos os grupos em questão. Nesse caso só haverá interseção quando tivermos passado:\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"name\": \"jose\",\n\"age\": \"30\",\n\"salary\": \"5001\"\n}\n` + "`" + `` + "`" + `` + "`" + `\nCaso seja colocado um valor maior que 5000 em *salary*, também será uma interseção. Temos como resultados esperados o valor de \"mygroup\": \"true\", \"taget_client\": \"true\" e   \"high_income\": \"true\".\n\n\n**Exemplo 6** **Regra com operações matemáticas**:\n\n- [Adição](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0012%20-%20increment_value/rules.featws): Nesse exemplo será somado um valor na variável. Com isso teremos como resposta de travel_1: travel_distance + 10, ou seja 10, e como resposta de travel_2: travel_distance + 10, ou seja, 110.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"travel_distance\": \"0\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n- Valor [Quadratico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0025%20-%20age%20_squared) e [Cubico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0026%20-%20age%20_cubed): Ao enviar um valor escolhido será retornado o valor cubico e quadrático do valor, podendo ser pedido um ou outro ao depender da regra escrita.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"idade\": \"5\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 7**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0034%20-%20simple_feature_with_condition) **Regra com condição**:\n\nNesta regra, o valor do recurso \"myfeat\" será retornado somente se a condição especificada for atendida. Se o corpo (body) enviado não cumprir essa condição, não será gerada nenhuma resposta em relação a esse recurso. Em outras palavras, a condição é um requisito para que a regra seja acionada e para que haja retorno do valor do recurso. Caso contrário, a regra não terá efeito e não será gerada nenhuma resposta em relação a ela.\n\nAo enviar um corpo (body) na requisição, se o valor de \"myothernumber\" for menor ou igual a 10, nenhum valor será retornado como resposta. Porém, se o valor for maior que 10, o valor de \"mynumber\" será retornado, acrescido de 12. Para obter uma resposta adequada, é necessário enviar ambos os valores no corpo da requisição. Isso garantirá que a condição seja atendida e a regra possa ser aplicada, gerando uma resposta com o valor desejado. É importante observar que o valor de \"myothernumber\" é um critério para a aplicação da regra e que o valor de \"mynumber\" é o resultado final que será retornado pela API.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"mynumber\": \"-5\",\n\"myothernumber\": \"12\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eval"
                ],
                "summary": "Evaluate the rulesheet / Avaliação da folha de Regra",
                "parameters": [
                    {
                        "description": "Parameters",
                        "name": "parameters",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.Eval"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/eval/{knowledgeBase}": {
            "post": {
                "security": [
                    {
                        "Authentication Api Key": []
                    }
                ],
                "description": "Para a realiza os testes basta clicar em *Try it out*, complete a folha de regra com os dados desejados e os demais campos caso necessário, em seguida, clique em *Execute*.\n\nA seguir, serão apresentados alguns exemplos de testes:\n\n[**Exemplo 1**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0001%20-%20one_parameter) **Testando uma variável**:\n` + "`" + `` + "`" + `` + "`" + `\n\"myboolfeat\" = \"$mynumber \u003c 12\"\n` + "`" + `` + "`" + `` + "`" + `\n\nNesse exemplo, é possível testar a feature *myboolfeat*. Ao abrir o arquivo [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0001%20-%20one_parameter/rules.featws), é possível observar que, se o valor de *mynumber* for menor que 12, a feature *myboolfeat* retornará *true*. Caso contrário, se for maior ou igual a 12, o retorno será *false*. Portanto, para testar essa regra, basta inserir o seguinte corpo de *Parameters*.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"mynumber\": \"1\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 2**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0003%20-%20simple_group) **Regra com um grupo**:\n\nNesse exemplo vamos testar a utilização de um grupo espeífico. Ao enviarmos um **clientingroup** que tenha no grupo [mygroup](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0003%20-%20simple_group/groups/mygroup.json), espera-se como retorno **mygroup** = true, caso seja passado uma agencia e conta válidas como a seguir:\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"branch\": \"00000\",\n\"account\": \"00000000\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 3**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0004%20-%20default_value_true) **Regra com um valor padrão**:\n\nNesse exemplo vamos testar a utilização de um valor padrão. Ao enviarmos um **gender** F sua resposta esperada será \"female\" = \"true\", e \"male\" =\"false\", pois no [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/rules.featws) da regra diz que female é o inverso de male. Logo caso seja enviado o valor de \"gender\": \"M\", teremos female =\"false\" e male =\"true\". Sendo enviado qualque outra letra como parâmetro será recebido female=\"false\" e male=\"true\", pois o default male está declarado primeiro do que o female no arquivo [features.json](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/features.json) da regra.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"gender\": \"F\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 4**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0005%20-%20precedence_order) **Regra com ordem de procedência**:\n\nNesse exemplo vamos testar a utilização de uma ordem de procedência. Similar ao caso anterior, se é enviado um valor menor ou igual à 18, receberemos como resposta \"menor_de_idade\": \"true\" e \"maior_de_idade\": \"false\", caso seja enviado um valor maior que 18 deveremos receber \"menor_de_idade\": \"false\" e \"maior_de_idade\": \"true\". Caso não seja enviado nenhum valor de idade, será interpretado como \"menor_de_idade\": \"true\" e \"maior_de_idade\":\"false\".\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"idade\": \"21\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 5**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0008%20-%20group%20_intersection) **Regra com interseção de grupos**:\n\nA interseção de grupos ocorre quando dois ou mais grupos têm elementos em comum, ou seja, há um conjunto de elementos que pertencem a todos os grupos em questão. Nesse caso só haverá interseção quando tivermos passado:\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"name\": \"jose\",\n\"age\": \"30\",\n\"salary\": \"5001\"\n}\n` + "`" + `` + "`" + `` + "`" + `\nCaso seja colocado um valor maior que 5000 em *salary*, também será uma interseção. Temos como resultados esperados o valor de \"mygroup\": \"true\", \"taget_client\": \"true\" e   \"high_income\": \"true\".\n\n\n**Exemplo 6** **Regra com operações matemáticas**:\n\n- [Adição](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0012%20-%20increment_value/rules.featws): Nesse exemplo será somado um valor na variável. Com isso teremos como resposta de travel_1: travel_distance + 10, ou seja 10, e como resposta de travel_2: travel_distance + 10, ou seja, 110.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"travel_distance\": \"0\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n- Valor [Quadratico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0025%20-%20age%20_squared) e [Cubico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0026%20-%20age%20_cubed): Ao enviar um valor escolhido será retornado o valor cubico e quadrático do valor, podendo ser pedido um ou outro ao depender da regra escrita.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"idade\": \"5\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 7**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0034%20-%20simple_feature_with_condition) **Regra com condição**:\n\nNesta regra, o valor do recurso \"myfeat\" será retornado somente se a condição especificada for atendida. Se o corpo (body) enviado não cumprir essa condição, não será gerada nenhuma resposta em relação a esse recurso. Em outras palavras, a condição é um requisito para que a regra seja acionada e para que haja retorno do valor do recurso. Caso contrário, a regra não terá efeito e não será gerada nenhuma resposta em relação a ela.\n\nAo enviar um corpo (body) na requisição, se o valor de \"myothernumber\" for menor ou igual a 10, nenhum valor será retornado como resposta. Porém, se o valor for maior que 10, o valor de \"mynumber\" será retornado, acrescido de 12. Para obter uma resposta adequada, é necessário enviar ambos os valores no corpo da requisição. Isso garantirá que a condição seja atendida e a regra possa ser aplicada, gerando uma resposta com o valor desejado. É importante observar que o valor de \"myothernumber\" é um critério para a aplicação da regra e que o valor de \"mynumber\" é o resultado final que será retornado pela API.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"mynumber\": \"-5\",\n\"myothernumber\": \"12\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eval"
                ],
                "summary": "Evaluate the rulesheet / Avaliação da folha de Regra",
                "parameters": [
                    {
                        "type": "string",
                        "description": "knowledgeBase",
                        "name": "knowledgeBase",
                        "in": "path"
                    },
                    {
                        "description": "Parameters",
                        "name": "parameters",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.Eval"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/eval/{knowledgeBase}/{version}": {
            "post": {
                "security": [
                    {
                        "Authentication Api Key": []
                    }
                ],
                "description": "Para a realiza os testes basta clicar em *Try it out*, complete a folha de regra com os dados desejados e os demais campos caso necessário, em seguida, clique em *Execute*.\n\nA seguir, serão apresentados alguns exemplos de testes:\n\n[**Exemplo 1**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0001%20-%20one_parameter) **Testando uma variável**:\n` + "`" + `` + "`" + `` + "`" + `\n\"myboolfeat\" = \"$mynumber \u003c 12\"\n` + "`" + `` + "`" + `` + "`" + `\n\nNesse exemplo, é possível testar a feature *myboolfeat*. Ao abrir o arquivo [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0001%20-%20one_parameter/rules.featws), é possível observar que, se o valor de *mynumber* for menor que 12, a feature *myboolfeat* retornará *true*. Caso contrário, se for maior ou igual a 12, o retorno será *false*. Portanto, para testar essa regra, basta inserir o seguinte corpo de *Parameters*.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"mynumber\": \"1\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 2**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0003%20-%20simple_group) **Regra com um grupo**:\n\nNesse exemplo vamos testar a utilização de um grupo espeífico. Ao enviarmos um **clientingroup** que tenha no grupo [mygroup](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0003%20-%20simple_group/groups/mygroup.json), espera-se como retorno **mygroup** = true, caso seja passado uma agencia e conta válidas como a seguir:\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"branch\": \"00000\",\n\"account\": \"00000000\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 3**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0004%20-%20default_value_true) **Regra com um valor padrão**:\n\nNesse exemplo vamos testar a utilização de um valor padrão. Ao enviarmos um **gender** F sua resposta esperada será \"female\" = \"true\", e \"male\" =\"false\", pois no [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/rules.featws) da regra diz que female é o inverso de male. Logo caso seja enviado o valor de \"gender\": \"M\", teremos female =\"false\" e male =\"true\". Sendo enviado qualque outra letra como parâmetro será recebido female=\"false\" e male=\"true\", pois o default male está declarado primeiro do que o female no arquivo [features.json](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/features.json) da regra.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"gender\": \"F\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 4**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0005%20-%20precedence_order) **Regra com ordem de procedência**:\n\nNesse exemplo vamos testar a utilização de uma ordem de procedência. Similar ao caso anterior, se é enviado um valor menor ou igual à 18, receberemos como resposta \"menor_de_idade\": \"true\" e \"maior_de_idade\": \"false\", caso seja enviado um valor maior que 18 deveremos receber \"menor_de_idade\": \"false\" e \"maior_de_idade\": \"true\". Caso não seja enviado nenhum valor de idade, será interpretado como \"menor_de_idade\": \"true\" e \"maior_de_idade\":\"false\".\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"idade\": \"21\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 5**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0008%20-%20group%20_intersection) **Regra com interseção de grupos**:\n\nA interseção de grupos ocorre quando dois ou mais grupos têm elementos em comum, ou seja, há um conjunto de elementos que pertencem a todos os grupos em questão. Nesse caso só haverá interseção quando tivermos passado:\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"name\": \"jose\",\n\"age\": \"30\",\n\"salary\": \"5001\"\n}\n` + "`" + `` + "`" + `` + "`" + `\nCaso seja colocado um valor maior que 5000 em *salary*, também será uma interseção. Temos como resultados esperados o valor de \"mygroup\": \"true\", \"taget_client\": \"true\" e   \"high_income\": \"true\".\n\n\n**Exemplo 6** **Regra com operações matemáticas**:\n\n- [Adição](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0012%20-%20increment_value/rules.featws): Nesse exemplo será somado um valor na variável. Com isso teremos como resposta de travel_1: travel_distance + 10, ou seja 10, e como resposta de travel_2: travel_distance + 10, ou seja, 110.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"travel_distance\": \"0\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n- Valor [Quadratico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0025%20-%20age%20_squared) e [Cubico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0026%20-%20age%20_cubed): Ao enviar um valor escolhido será retornado o valor cubico e quadrático do valor, podendo ser pedido um ou outro ao depender da regra escrita.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"idade\": \"5\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n\n\n[**Exemplo 7**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0034%20-%20simple_feature_with_condition) **Regra com condição**:\n\nNesta regra, o valor do recurso \"myfeat\" será retornado somente se a condição especificada for atendida. Se o corpo (body) enviado não cumprir essa condição, não será gerada nenhuma resposta em relação a esse recurso. Em outras palavras, a condição é um requisito para que a regra seja acionada e para que haja retorno do valor do recurso. Caso contrário, a regra não terá efeito e não será gerada nenhuma resposta em relação a ela.\n\nAo enviar um corpo (body) na requisição, se o valor de \"myothernumber\" for menor ou igual a 10, nenhum valor será retornado como resposta. Porém, se o valor for maior que 10, o valor de \"mynumber\" será retornado, acrescido de 12. Para obter uma resposta adequada, é necessário enviar ambos os valores no corpo da requisição. Isso garantirá que a condição seja atendida e a regra possa ser aplicada, gerando uma resposta com o valor desejado. É importante observar que o valor de \"myothernumber\" é um critério para a aplicação da regra e que o valor de \"mynumber\" é o resultado final que será retornado pela API.\n` + "`" + `` + "`" + `` + "`" + `\n{\n\"mynumber\": \"-5\",\n\"myothernumber\": \"12\"\n}\n` + "`" + `` + "`" + `` + "`" + `\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eval"
                ],
                "summary": "Evaluate the rulesheet / Avaliação da folha de Regra",
                "parameters": [
                    {
                        "type": "string",
                        "description": "knowledgeBase",
                        "name": "knowledgeBase",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "version",
                        "name": "version",
                        "in": "path"
                    },
                    {
                        "description": "Parameters",
                        "name": "parameters",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.Eval"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.Eval": {
            "type": "object",
            "additionalProperties": true
        }
    },
    "securityDefinitions": {
        "Authentication Api Key": {
            "type": "apiKey",
            "name": "X-API-Key",
            "in": "header"
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "FeatWS Ruler",
	Description:      "O projeto Ruler é uma implementação do motor de regras [grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine), que é utilizado para avaliar regras no formato .grl . O Ruler permite que as regras definidas em arquivos .grl sejam avaliadas de maneira automática e eficiente, ajudando a automatizar as decisões tomadas pelo FeatWS. Isso possibilita que o sistema possa analisar e classificar grandes quantidades de informações de maneira rápida e precisa.\n\nAo utilizar as regras fornecidas pelo projeto Ruler, o FeatWS é capaz de realizar análises de regras em larga escala e fornecer resultados precisos e relevantes para seus usuários. Isso é especialmente importante em áreas como análise de sentimentos em mídias sociais, detecção de fraudes financeiras e análise de dados em geral.\n\nAntes de realizar os testes no Swagger, é necessário autorizar o acesso clicando no botão **Authorize**, ao lado, e inserindo a senha correspondente. Após inserir o campo **value** e clicar no botão **Authorize**, o Swagger estará disponível para ser utilizado.\n\nA seguir é explicado com mais detalhes sobre os endpoints:\n- **/Eval**: Esse endpoint é utilizado apenas para aplicações que possuem uma única folha de regra padrão.\n- **/Eval/{knowledgeBase}**: Nesse endpoint, é necessário informar o parâmetro com o nome da folha de regra desejada e, como resultado, será retornado a última versão da folha de regra correspondente.\n- **/Eval/{knowledgeBase}/{version}**: Nesse endpoint é necessário colocar o parâmetro do nome da folha de regra como também o número da versão da folha de regra que você deseja testar a regra.\n\n**Parameters / Parâmetros**\nNo **knowledgeBase**, você pode especificar o nome da folha de regras que deseja utilizar. Já o **version** você coloca a versão que você deseja avaliar. Em **Paramenter**, é possível enviar os parametros que você deseja testar na folha de regra.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
