start:
		helm repo add kong https://charts.konghq.com
		helm repo update
		helm install kong kong/kong \
			--namespace aerona \
			--set ingressController.installCRDs=false \
			--set postgresql.enabled=true \
			--set env.database=postgres \
			--set env.pg_user=kong \
			--set env.pg_password=kong \
			--set env.pg_database=kong \
			--set proxy.type=LoadBalancer
