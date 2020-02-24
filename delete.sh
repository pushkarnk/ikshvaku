kubectl delete -f deploy/crds/cache.example.com_v1alpha1_ikshvaku_cr.yaml
kubectl delete -f deploy/operator.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/crds/cache.example.com_ikshvakus_crd.yaml
