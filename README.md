# Middleware Publish/Subscribe (Pub/Sub)

## Objetivo do Projeto

Este projeto tem como objetivo implementar um middleware do tipo Publish/Subscribe (Pub/Sub) para a disciplina de Programação Distribuída. O sistema permite a comunicação assíncrona entre publicadores e assinantes através de tópicos, utilizando sockets TCP e protocolo JSON.

O middleware gerencia tópicos dinamicamente, realiza bufferização de mensagens, encaminha publicações para todos os assinantes de cada tópico e descarta mensagens sem assinantes, informando o publicador quando isso ocorre. O foco é simplicidade, clareza, concorrência correta e aderência ao enunciado da atividade.

## O que o projeto faz
- Permite múltiplas conexões TCP simultâneas
- Gerencia tópicos de forma dinâmica
- Realiza operações de publish, subscribe e unsubscribe
- Encaminha mensagens para todos os assinantes do tópico
- Descarta mensagens sem assinantes e informa o publicador
- Bufferiza mensagens para cada tópico
- Utiliza goroutines, channels e mutexes para concorrência
- Comunicação via protocolo JSON newline-delimited

## Informações
- **Disciplina:** Programação Distribuída 2026.1
- **Professor:** Ruan Delgado Gomes
- **Alunas:** Joana Elise e Maria Eduarda Vitorino
