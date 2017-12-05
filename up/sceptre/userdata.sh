#!/bin/sh
version={version}
set -o xtrace
# install docker-ce
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
apt-get update
apt-get -y install docker-ce
# cw-agent.service-start
cat >> /etc/systemd/system/cw-agent.service << EOF
[Install]
WantedBy=multi-user.target
[Unit]
Description=cloudwatch agent
After=docker.service
Requires=docker.service
[Service]
Restart=always
ExecStartPre=-/usr/bin/docker kill cw-agent
ExecStartPre=/usr/bin/docker pull cscr/cw-agent
ExecStart=/usr/bin/docker run -i --rm --name cw-agent -v /var/log:/var/log cscr/cw-agent
ExecStop=/usr/bin/docker stop cw-agent
EOF
# cw-agent.service-end
EcrUrl={ecr_domain}/{ecr_tag}:{version}
# deployed.service-start
cat >> /etc/systemd/system/deployed.service << EOF
[Install]
WantedBy=multi-user.target
[Unit]
Description=deployed docker service
After=docker.service
Requires=docker.service
[Service]
TimeoutStartSec=0
Restart=always
ExecStartPre=-/usr/bin/docker kill deployed
ExecStartPre=/usr/bin/docker pull mesosphere/aws-cli
ExecStartPre=/bin/sh -c 'eval \$(docker run -i mesosphere/aws-cli ecr get-login --no-include-email --region {region})'
ExecStartPre=/usr/bin/docker pull $EcrUrl
ExecStart=/usr/bin/docker run -i --rm --name deployed -p 80:8080 $EcrUrl
ExecStop=/usr/bin/docker stop deployed
EOF
# deployed.service-end
# make sure .service is loaded, enable and start
systemctl daemon-reload
systemctl --no-block enable deployed.service cw-agent.service
systemctl --no-block start deployed.service cw-agent.service
