# Nullable Types

## Proposta
Construir e consumir [APIs REST](https://www.redhat.com/en/topics/api/what-is-a-rest-api) são tarefas comuns no dia a dia de um desenvolvedor. Porém, nem sempre os dados retornados por essas APIs são completos, podendo conter valores nulos. Em GO, os tipos básicos não aceitam valores nulos, o que pode ser um desafio ao lidar com APIs REST.

A seguir descrevo em detalhes a problemática e duas abordagens para solucionar.

## Requisitos
Vamos ao cenario:
```yaml
Vamos simular uma API que:
 1. recebe um requisição POST;
    1.1 não é necessário validar o verbo HTTP;
 2. obtem o parametro 'count' do tipo inteiro do payload;
    2.1. caso o payload esteja mal formatado, retorna um erro 400;
		2.2. caso o cliente não informe o campo 'count', retorna um erro 422;
 3. retorna como resposta o valor 'count' acrescido de 1.
```

### Testes de aceitação
Os [cenarios de teste](./acceptance/types_test.go#22) descrevem o comportamento esperado da API, conforme os requisitos acima.

## 1. Back to Basics
A [primeira abordagem](./types/basic.go#12) é a mais simples e direta. Utilizamos uma [struct](./types/basic.go#8) contendo o atributo `Count` do tipo `int` para parsear o payload da requisição. Caso o parse falhe, retornamos um erro 400. Na sequencia incrementamos o valor de `Count` e retornamos o resultado.
```go
type PayloadWithBasicTypes struct {
	Count int `json:"count"`
}
```

### 1.1. Testes
```bash
$ go test -run ^TestAcceptance$/^with_basic_type$ nullable/acceptance
```
Dois testes falharam. Ambos esperavam `422` e obtiveram `200`, o primeiro não informando o campo `count`e o segundo informando como `null`. Ao mapearmos o conteudo do payload para a struct, os campos não informados são inicializados com [zero values](https://go.dev/tour/basics/12), logo incrementamos o valor `0` em 1 e retornamos 200, não atendendo aos cenarios em que o campo `count` não é informado.

Mas como identificar se o cliente deixou de informar o campo ou se simplesmente informou com valor `0`?

## 2. Pointing to the solution
Precisamos então de um tipo que mantenha a informação (int) do campo `count` e um indicativo de que o campo foi (ou não) preenchido. A [segunda abordagem](./types/pointer.go#12) utiliza um ponteiro para `int` para atender a essa necessidade.

```go
type PayloadWithPointerTypes struct {
	Count *int `json:"count"`
}
```
Com esse novo payload, muito pouco se altera na implementação, adicionamos apenas um [if](./types/pointer.go#19) para verificar se o campo `count` foi preenchido e retornando erro 422 caso contrário.
```go
  if payload.Count == nil {
    http.Error(writer, "", http.StatusUnprocessableEntity)
    return
  }
```

### 2.1. Testes
```bash
$ go test -run ^TestAcceptance$/^with_pointer_type$ nullable/acceptance
```
Nossos testes de aceitação agora passam. :white_check_mark: :tada: :sparkles:

Agora podemos ir tomar aquele café...

![alt](../static/ou_sera_que_nao.jpg)

### 3. Customizando seu ponteiro
A abordagem descrita anteriormente e suficiente para a maioria dos casos, mas gostaria de listar algums pontos de atenção:
  - Adicionar a possibilidade de nulo ao campo 'count' obriga o desenvolvedor a verificar se o campo e nulo a cada leitura e potencialmente antes de escrever também.
  - Conforme a estrutura Payload cresce, o custo de verificar todos os campos nullables torna-se dificil de gerencial e propenso a erros.
  - A estrutura Payload pode ser usada em varios lugares, o que significa que a verificacao de nulos deve ser feita em varios lugares.

Então podemos ir além, com uma abordagem mais "limpa" e criar um tipo personalizado que encapsula o `*int` e fornece metodos para manipula-lo.
Nossa [terceira abordagem](./types/custom.go#40) utiliza um tipo `NullableInt` como substituto para `*int`.
```go
type NullableInt struct {
	value *int
}
type PayloadWithCustomTypes struct {
	Count NullableInt `json:"count"`
}
```

Mantemos o valor `value` do ponteiro privado e fornecemos metodos para manipula-lo.
```go
// Obtem o valor, se nulo retorna 0, evitando nil pointer dereference
func (t NullableInt) Value() int
// Permite atualizar o valor. Não foi adicionado a verficação de nulo, mas em situações reais seria necessario.
func (t *NullableInt) Set(v int)
```

Temos ainda a implementação dos metodos das interfacec [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler) e [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) para permitir a [serialização](./types/custom.go#31) e [deserialização](./types/custom.go#27) do tipo `NullableInt` como um `*int`.

### 3.1. Testes
```bash
$ go test -run ^TestAcceptance$/^with_custom_type$ nullable/acceptance
```
Os testes seguem passando, mas agora temos uma estrutura mais robusta e facil de manter. :white_check_mark: :tada: :sparkles:
Vale ressaltar que existem pacotes que implementam tipo com comportamentos similares, como [sql.NullInt64](https://pkg.go.dev/database/sql#NullInt64) utilizado para representar valores nulos em banco de dados e [json.Number](https://pkg.go.dev/encoding/json#Number) que, assim como implementamos, representa valores numericos em JSON, permitindo valores nulos e com diversas operações de manipulação.
