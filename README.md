# sms-prometheus
An app for send sms when trigger an alert in prometheus
with this application you can send sms to a specific URL when one of the rules triggered in prometheus.
after trigger a rule in prometheus ,prometheus send a warning to alertmanager with alert title and some other usefull datas.
with that in comming request from prometheus ,alertmanager route that config and in last and most important step we read request that alertmanager sent to our app and then send alert to users tel number.
