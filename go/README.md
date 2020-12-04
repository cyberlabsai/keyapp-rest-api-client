# Golang
Esta aplicação exemplo, mostra como fazer a comunicação com a API KeyApp da Cyberlabs.

## Read the docs!
[Neste link](https://api.keyapp.ai/docs) você tem acesso a documentação oficial da API.

## Open API & Swagger
[Neste link](https://api.keyapp.ai/) você tem a referência dos enpoints disponíveis e é possível testá-los.

## Exemplos disponíveis
Esta aplicação demonstra como validar um token de uma requisição assinada feita quando uma ação acontece no KeyApp. Para criar uma ação assinada, é preciso ter uma conta no [KeyApp](https://play.google.com/store/apps/details?id=ai.cyberlabs.keyapp) e possuir um portal.

Além disso, é preciso [criar uma aplicação](https://api.keyapp.ai/docs/autenticacao/#geracao-das-credenciais) que dará acesso as suas credenciais de aplicação (chave de API e o segredo).

### Em que momento recebo o token na requisição?
No cadastro de sua ação assinada, uma URL de destino é fornecida. Ou seja, quando alguém executar uma ação do Portal, seja para abrir uma porta ou acenda uma luz, uma requisição HTTP irá chegar na URL cadastrada na ação. Nessa requisição, em uma header *Authorization* um token será fornecido.

O token fornecido na requisição da ação é uma prova que essa requisição foi feita pelos servidores do KeyApp. Com esse token, você pode solicitar o *id* do usuário que executou a ação e o *id* do Portal de qual a ação pertence. Além disso, se o token não foi gerado pelos servidores do KeyApp, um resposta de token inválido será retornada.