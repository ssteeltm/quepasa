#!/bin/bash

cd /root/.n8n
/usr/bin/n8n import:workflow --input=/opt/quepasa-source/extra/n8n+chatwoot/ --separate
/usr/bin/n8n update:workflow --id 1008 --active=true
/usr/bin/n8n update:workflow --id 1009 --active=true
/usr/bin/n8n update:workflow --id 1010 --active=true
/usr/bin/n8n update:workflow --id 1011 --active=true

exit 0