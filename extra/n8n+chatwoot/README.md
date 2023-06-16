[![Chatwoot Profile Update]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/ChatwootProfileUpdate.json)
[![Chatwoot To Quepasa]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/ChatwootToQuepasa.json)
[![Chatwoot To Quepasa Greetings]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/ChatwootToQuepasaGreetings.json)
[![Post To Chatwoot]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/PostToChatwoot.json)
[![Quepasa Automatic]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaAutomatic.json)
[![Quepasa Chat Control]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaChatControl.json)
[![Quepasa Inbox Control]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaInboxControl.json)
[![Quepasa Qrcode]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaQrcode.json)
[![Quepasa To Chatwoot]](https://raw.githubusercontent.com/nocodeleaks/quepasa/main/extra/n8n%2Bchatwoot/QuepasaToChatwoot.json)

Use N8N Environment File to set these variables:
(use your respectives ids)

# enable a single conversation per contact, for all life time, not just a ticket
C8Q_SINGLETHREAD=false

# used by ChatwootToQuepasa
C8Q_QUEPASAINBOXCONTROL=8

# used by ChatwootToQuepasa | QuepasaToChatwoot
C8Q_QUEPASACHATCONTROL=7

# used by ChatwootToQuepasa
C8Q_CHATWOOTPROFILEUPDATE=5

# used by QuepasaToChatwoot
C8Q_POSTTOCHATWOOT=6

# used by QuepasaToChatwoot
C8Q_POSTTOWEBCALLBACK=11

# used by QuepasaToChatwoot
C8Q_CHATWOOTTOQUEPASAGREETINGS=4

C8Q_QP_DEFAULT_USER=

# get initial from $N8N_HOST, otherwise from here
C8Q_CW_PUBLIC_URL=dominiochatwoot

# (optional) set title for sincrony bot for messages from whatsapp
# C8Q_QP_BOTTITLE=

# (optional) set identifier for quepasa inbox control contact
# C8Q_QP_CONTACT=