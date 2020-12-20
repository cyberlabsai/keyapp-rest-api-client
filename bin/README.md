# Scripts

## get-access-token.sh

Esse script funciona como uma espécie de *wizard* e executa todos os
passos necessários para conseguir-se um *Access Token*.

Você pode "pular" alguns passos disponibilizando certas variáveis de
ambiente. Por exemplo: caso você já tenha um *typecode*, pode simplesmente
informá-lo diretamente, como em

```bash
$ TYPECODE=o-typecode-que-você-já-tem bin/get-access-token.sh
```

Dessa forma você não precisará receber um novo e-mail com seu código.


O mesmo funciona para a variável `TICKET`.


## create-app-key.sh

Uma vez que você tenha seu *Access Token* em mãos, esse script serve para
que você possa gerar uma nova Chave de API a ser usada em alguma aplicação
sua.


`Usage: bin/create-app-key.sh access_token key_name`

`access_token` é o seu *Access Token* e `key_name` é o nome da Chave que
você quer gerar (geralmente as pessoas dão à Chave de API o mesmo nome que
dão à aplicação que a utiliza).
