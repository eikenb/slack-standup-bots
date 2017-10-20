StackName := up-bot-redis

redis-create: _call-stack
redis-create: stack-call = create-stack
redis-create: wait-call = stack-create-complete

redis-update: _call-stack
redis-update: stack-call = update-stack
redis-update: wait-call = stack-update-complete

_call-stack:
	aws cloudformation ${stack-call} --stack-name ${StackName} \
		--tags "$$(cat tags.json)" \
		--parameters \
			ParameterKey=VpcId,ParameterValue="$$(cf.sh vpc)"\
			ParameterKey=SubnetIds,ParameterValue="'$$(cf.sh subs resource)'"\
		--template-body "file://redis-cf.yaml"
	aws cloudformation wait ${wait-call} --stack-name ${StackName}

