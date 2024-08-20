### Workflows links

* [(1.0.1) Chatwoot Profile Update](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/ChatwootProfileUpdate.json)
* [(1.0.32) Chatwoot To Quepasa](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/ChatwootToQuepasa.json)
* [(1.0.17) Post To Chatwoot](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/PostToChatwoot.json)
* [(1.0.6) Quepasa Automatic](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaAutomatic.json)
* [(1.0.6) Quepasa Chat Control](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaChatControl.json)
* [(1.0.9) Quepasa Inbox Control](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaInboxControl.json)
* [(1.0.1) Quepasa Qrcode](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaQrcode.json)
* [(1.0.18) Quepasa To Chatwoot](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaToChatwoot.json)
* [(1.0.13) Get Chatwoot Contacts](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/GetChatwootContacts.json)
* [(1.0.1) Post To WebCallBack](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/PostToWebCallBack.json)
* [(1.0.1) Chatwoot Extra](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/ChatwootExtra.json)
* [(1.0.0) To Chatwoot Transcript]()


### Use N8N Environment File to set these variables:
> use your respectives ids

	# C8Q_QUEPASAINBOXCONTROL
	> (integer) Workflow Id - QuepasaInboxControl, used by ChatwootToQuepasa
	> (default) 1001

	# C8Q_GETCHATWOOTCONTACTS
	> (integer) Workflow Id - GetChatwootContacts, (2023/06/19) used by QuepasaToChatwoot
	> (default) 1002
		
	# C8Q_QUEPASACHATCONTROL
	> (integer) Workflow Id - QuepasaChatControl, used by ChatwootToQuepasa | QuepasaToChatwoot
	> (default) 1003
	
	# C8Q_CHATWOOTPROFILEUPDATE
	> (integer) Workflow Id - ChatwootProfileUpdate, used by ChatwootToQuepasa
	> (default) 1004
	
	# C8Q_POSTTOWEBCALLBACK
	> (integer) Workflow Id - PostToWebCallBack, used by QuepasaToChatwoot
	> (default) 1005
	
	# C8Q_POSTTOCHATWOOT
	> (integer) Workflow Id - PostToChatwoot, used by QuepasaToChatwoot
	> (default) 1006
	
	# C8Q_CHATWOOTTOQUEPASAGREETINGS
	> (integer) Workflow Id - ChatwootToQuepasaGreetings, used by QuepasaToChatwoot
	> (default) 1007
	
	# C8Q_CHATWOOTEXTRA
	> (string) Workflow Id - ChatwootExtra, used by QuepasaQrcode
	> (default) iiEsUj7ybtzEZAFj
	
	# C8Q_TOCHATWOOTTRANSCRIPT
	> (string) Workflow Id - ToChatwootTranscript, used by transcript audios to text
	> (default) pi4APHD9F05Dv6FR
	
	# C8Q_TOCHATWOOTTRANSCRIPTRESUME
	> (boolean => true | false) Gets a resume (OpenAI) for audio transcripted messages

	# C8Q_SINGLETHREAD
	> (boolean => true | false) Enable a single conversation per contact, for all life time, not just a ticket
	
	# C8Q_QP_DEFAULT_USER
	> (optional) (string => user@domain.com)

	# C8Q_CW_PUBLIC_URL
	> (optional) (string => https://chatwoot.com) get initial from $N8N_HOST, otherwise from here, used by ChatwootToQuepasa

	# C8Q_QP_BOTTITLE
	> (optional) (string => QuePasa [the best]) set title for sincrony bot for messages from whatsapp

	# C8Q_QP_CONTACT
	> (optional) (string => control@quepasa.io) set identifier for quepasa inbox control contact
	
	# C8Q_SUFFICIT_CONTEXTID
	> (optional) (string => 00000000-0000-0000-0000-000000000000) sufficit client identification

	# C8Q_UNKNOWN_SENDER
	> (optional) 'Mensagem de sistema: '
	
	# C8Q_MSGFOR_UNKNOWN_CONTENT
	> (optional) '"Algum EMOJI" ou "Alguma Rea&#231;&#227;o que o sistema n&#227;o entende ainda ..."'
	
	# C8Q_MSGFOR_EDITED_CONTENT
	> (optional) '*** Mensagem editada ***'
	
	# C8Q_MSGFOR_ATTACHERROR_CONTENT
	> (optional) 'Falha ao baixar anexo'
	
	# C8Q_MSGFOR_LOCALIZATION_CONTENT
	> (optional) 'Envio de localiza&#231;&#227;o'
	
	# C8Q_MSGFOR_REVOKED_CONTENT
	> (optional) '*** Mensagem apagada ***'
	
	# C8Q_MSGFOR_CALL_CONTENT
	> (optional) 'Usu&#225;rio requisitou uma chamada de voz | video'
	
	# C8Q_MSGFOR_NO_CSAT
	> (optional) 'Atendimento concluído'

### Use ChatWoot INBOX Webhook parameters:
> individual inboxes

	# pat
	> (boolean) Prepend Agent Title - default true - should include agent title before msg content ?
	
	# st | singlethread
	> (boolean) Single Thread - default false - should use a single thread conversation and do not open a new ?