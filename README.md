# sls-rtc-backend


[こちら](https://hogehoge-banana.xyz/playground/sls-rtc/)のバックエンドの実装です。
AWS SAMを利用させていただいております。
とにかくwebrtcを試してみたいために作ったのでいろいろザルです。
ぶっちゃけよく見ればいろいろバレて、いろいろできてしまいますが、
完全に趣味で私のポケットマネーで運営してるのでどうかいじめないでください。お願いします。

## table design

dynamodbのレコードイメージ


参加者: participant
ws接続: connection
meetスペース: room

aaa: connect

|   pk              | wpush-p256dh | wpush-auth | participantID |
|-------------------| -------------|----------| ------------- |
| participantID:aaa | asd          | qwer     |aaa           |

issue participantId from backend: aaa

aaa: enter room

|   pk              | wpush-p256dh | wpush-auth | participantID | roomId |
|-------------------| -------------|----------| ------------- | ------ |
| participantID:aaa | asd          | qwer     |aaa            | room1 |


bbb: connect

|   pk              | wpush-p256dh | wpush-auth | participantID | roomId |
|-------------------| -------------|----------| ------------- | ------ |
| participantID:aaa | asd          | qwer     |aaa            | room1  |
| participantID:bbb | hjk          | yuio     |bbb            |        |


bbb: enter room

|   pk              | wpush-p256dh | wpush-auth | participantID | roomId |
|-------------------| -------------|----------| ------------- | ------ |
| participantID:aaa | asd          | qwer     |aaa            | room1  |
| participantID:bbb | hjk          | yuio     |bbb            | room1  |



GSI

| roomID(KeySchema) | connectionID(projection)  | userName(projection)   |
|-------------------|---------------------------|----------|
| johnroom | 78iujhyt542qw | john |
| johnroom | 4edfgtredf0ol | michel |


## 使い方

`wscat` コマンドを使ったデバッグの例です

### connect

接続

```
wscat -c wss://{{ api gateway endpoint }}/slsrtc
```

接続に成功すると接続IDとともに通知を受信します

```
{"type":"connected","body":"{your connection id}"}
```

### create room

某whereなんとかとかいうビデオ会議のapiではrestで提供されているあれです。
課金とかセキュリティとか、レスポンスのハンドリング考えるとrestのほうがいいかもしれないです。
~~めんどくさいのでwebsocketで全部実装してしまいました。~~

```
{ "action": "createroom" }
```

ルームIDとともに`room-created` イベントを受信するのでそのIDを控えてください

```
{"type":"room-created","body":"{generated room id}"}
```

### enter room

create room で作成したroom id を指定してroomに参加します。

```
{"action":"enterroom","roomID":"{generated room id}"}
```

参加に成功するとすべてのroomメンバーに`enter`イベントがブロードキャストされます

```
{"type":"enter","roomID":"{room id}","connectionID":"{your connection id}"}
```

このイベントをトリガーにWebRTCのシグナリングを開始します。

### leave room

誰かが接続が切れると`leave`イベントが同じルームのコネクションにブロードキャストされます。


```
{"type":"leave","roomID":"{room id}","connectionID":"{your connection id}"}
```

このイベントを利用してWebRTCのコネクションを破棄します。

## Reference

https://aws.amazon.com/jp/blogs/news/simulating-amazon-dynamodb-unique-constraints-using-transactions/


