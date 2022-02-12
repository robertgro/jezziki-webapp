# jezziki-webapp

discontinued but finished web app utilizing a golang backend and reactjs frontend.

# dependencies

react 17.0.1  
golang 1.16  
postgresql 13.1  

# packages

echo/v4 (web server node provider and proxy using a controller)  
casbin/v2 (rule model based basic auth)  
ntp (ntp server time sync)  
pq (postgres sql driver package)  
ksuid (global unique identifier generation)  
crypto/rand + math/big (token generation)  
yaml.v2 (for app-config.yaml parsing)  
graceful (graceful shutdown)  
net/sys/rate (limiter)/validator0.9/encoding json/reflect (db type and object reflection)/regexp/strconv/  

admin dash control uses a MutationObserver object, is written in vanilla javascript and includes Jodit v.3 an open source WYSIWYG html editor  

# additional content

files to create CentOS8 service workers  
unix shellscripts to dump logs and sql database  
there is a template renderer engine for .gohtml template files and the server has proper error handling  
the proxy mirroring multiple app instances and collects unique visitor info by querying via whois tool  
