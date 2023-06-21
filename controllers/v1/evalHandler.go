package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// EvalHandler godoc
// @Summary 		Evaluate the rulesheet / Avaliação da folha de Regra
// @Description     Para a realiza os testes basta clicar em *Try it out*, complete a folha de regra com os dados desejados e os demais campos caso necessário, em seguida, clique em *Execute*.
// @Description
// @Description		A seguir, serão apresentados alguns exemplos de testes:
// @Description
// @Description		[**Exemplo 1**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0001%20-%20one_parameter) **Testando uma variável**:
// @Description		```
// @Description			"myboolfeat" = "$mynumber < 12"
// @Description		```
// @Description
// @Description		Nesse exemplo, é possível testar a feature *myboolfeat*. Ao abrir o arquivo [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0001%20-%20one_parameter/rules.featws), é possível observar que, se o valor de *mynumber* for menor que 12, a feature *myboolfeat* retornará *true*. Caso contrário, se for maior ou igual a 12, o retorno será *false*. Portanto, para testar essa regra, basta inserir o seguinte corpo de *Parameters*.
// @Description		```
// @Description		{
// @Description			"mynumber": "1"
// @Description		}
// @Description		```
// @Description
// @Description
// @Description		[**Exemplo 2**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0003%20-%20simple_group) **Regra com um grupo**:
// @Description
// @Description		Nesse exemplo vamos testar a utilização de um grupo espeífico. Ao enviarmos um **clientingroup** que tenha no grupo [mygroup](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0003%20-%20simple_group/groups/mygroup.json), espera-se como retorno **mygroup** = true, caso seja passado uma agencia e conta válidas como a seguir:
// @Description		```
// @Description		{
// @Description			"branch": "00000",
// @Description	 		"account": "00000000"
// @Description		}
// @Description		```
// @Description
// @Description
// @Description		[**Exemplo 3**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0004%20-%20default_value_true) **Regra com um valor padrão**:
// @Description
// @Description		Nesse exemplo vamos testar a utilização de um valor padrão. Ao enviarmos um **gender** F sua resposta esperada será "female" = "true", e "male" ="false", pois no [rules.featws](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/rules.featws) da regra diz que female é o inverso de male. Logo caso seja enviado o valor de "gender": "M", teremos female ="false" e male ="true". Sendo enviado qualque outra letra como parâmetro será recebido female="false" e male="true", pois o default male está declarado primeiro do que o female no arquivo [features.json](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0004%20-%20default_value_true/features.json) da regra.
// @Description		```
// @Description		{
// @Description			"gender": "F"
// @Description		}
// @Description		```
// @Description
// @Description
// @Description		[**Exemplo 4**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0005%20-%20precedence_order) **Regra com ordem de procedência**:
// @Description
// @Description		Nesse exemplo vamos testar a utilização de uma ordem de procedência. Similar ao caso anterior, se é enviado um valor menor ou igual à 18, receberemos como resposta "menor_de_idade": "true" e "maior_de_idade": "false", caso seja enviado um valor maior que 18 deveremos receber "menor_de_idade": "false" e "maior_de_idade": "true". Caso não seja enviado nenhum valor de idade, será interpretado como "menor_de_idade": "true" e "maior_de_idade":"false".
// @Description		```
// @Description		{
// @Description			"idade": "21"
// @Description		}
// @Description		```
// @Description
// @Description
// @Description		[**Exemplo 5**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0008%20-%20group%20_intersection) **Regra com interseção de grupos**:
// @Description
// @Description		A interseção de grupos ocorre quando dois ou mais grupos têm elementos em comum, ou seja, há um conjunto de elementos que pertencem a todos os grupos em questão. Nesse caso só haverá interseção quando tivermos passado:
// @Description		```
// @Description		{
// @Description			"name": "jose",
// @Description			"age": "30",
// @Description			"salary": "5001"
// @Description		}
// @Description		```
// @Description		Caso seja colocado um valor maior que 5000 em *salary*, também será uma interseção. Temos como resultados esperados o valor de "mygroup": "true", "taget_client": "true" e   "high_income": "true".
// @Description
// @Description
// @Description		**Exemplo 6** **Regra com operações matemáticas**:
// @Description
// @Description		- [Adição](https://github.com/bancodobrasil/featws-transpiler/blob/develop/__tests__/cases/0012%20-%20increment_value/rules.featws): Nesse exemplo será somado um valor na variável. Com isso teremos como resposta de travel_1: travel_distance + 10, ou seja 10, e como resposta de travel_2: travel_distance + 10, ou seja, 110.
// @Description		```
// @Description		{
// @Description			"travel_distance": "0"
// @Description		}
// @Description		```
// @Description		- Valor [Quadratico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0025%20-%20age%20_squared) e [Cubico](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0026%20-%20age%20_cubed): Ao enviar um valor escolhido será retornado o valor cubico e quadrático do valor, podendo ser pedido um ou outro ao depender da regra escrita.
// @Description		```
// @Description		{
// @Description			"idade": "5"
// @Description		}
// @Description		```
// @Description
// @Description
// @Description		[**Exemplo 7**](https://github.com/bancodobrasil/featws-transpiler/tree/develop/__tests__/cases/0034%20-%20simple_feature_with_condition) **Regra com condição**:
// @Description
// @Description		Nesta regra, o valor do recurso "myfeat" será retornado somente se a condição especificada for atendida. Se o corpo (body) enviado não cumprir essa condição, não será gerada nenhuma resposta em relação a esse recurso. Em outras palavras, a condição é um requisito para que a regra seja acionada e para que haja retorno do valor do recurso. Caso contrário, a regra não terá efeito e não será gerada nenhuma resposta em relação a ela.
// @Description
// @Description		Ao enviar um corpo (body) na requisição, se o valor de "myothernumber" for menor ou igual a 10, nenhum valor será retornado como resposta. Porém, se o valor for maior que 10, o valor de "mynumber" será retornado, acrescido de 12. Para obter uma resposta adequada, é necessário enviar ambos os valores no corpo da requisição. Isso garantirá que a condição seja atendida e a regra possa ser aplicada, gerando uma resposta com o valor desejado. É importante observar que o valor de "myothernumber" é um critério para a aplicação da regra e que o valor de "mynumber" é o resultado final que será retornado pela API.
// @Description		```
// @Description		{
// @Description			"mynumber": "-5",
// @Description			"myothernumber": "12"
// @Description		}
// @Description		```
// @Description
// @Tags 			eval
// @Accept  		json
// @Produce  		json
// @Param			knowledgeBase path string false "knowledgeBase"
// @Param 			version path string false "version"
// @Param  			parameters body payloads.Eval true "Parameters"
// @Success 		200 {string} string "ok"
// @Failure 		400,404 {object} string
// @Failure 		500 {object} string
// @Failure 		default {object} string
// @Security 		Authentication Api Key
// @Router 			/eval/{knowledgeBase}/{version} [post]
// @Router 			/eval/{knowledgeBase} [post]
// @Router 			/eval [post]
// This function handles requests to evaluate a knowledge base and returns the result as a JSON object.
func EvalHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		knowledgeBaseName := c.Param("knowledgeBase")
		if knowledgeBaseName == "" {
			knowledgeBaseName = services.DefaultKnowledgeBaseName
		}

		version := c.Param("version")
		if version == "" {
			version = services.DefaultKnowledgeBaseVersion
		}

		log.Debugf("Eval with %s %s\n", knowledgeBaseName, version)

		knowledgeBase, requestError := services.EvalService.GetKnowledgeBase(knowledgeBaseName, version)
		if requestError != nil {
			c.String(requestError.StatusCode, requestError.Message)
			return
		}

		decoder := json.NewDecoder(c.Request.Body)
		var t payloads.Eval
		err := decoder.Decode(&t)
		if err != nil {
			log.Errorf("Erro on json decode: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprint(c.Writer, "Error on json decode")
			return
		}
		log.Debugln(t)

		ctx := types.NewContextFromMap(t)
		ctx.RawContext = c.Request.Context()

		result, err := services.EvalService.Eval(ctx, knowledgeBase)
		if err != nil {

			log.Errorf("Error on eval: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprint(c.Writer, "Error on eval")
			return
		}

		log.Debug("Context:\n\t", ctx.GetEntries(), "\n\n")
		log.Debug("Features:\n\t", result.GetFeatures(), "\n\n")

		responseCode := http.StatusOK

		if result.Has("requiredParamErrors") {
			responseCode = http.StatusBadRequest
		}

		c.JSON(responseCode, result.GetFeatures())
	}

}
