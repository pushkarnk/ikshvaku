# ikshvaku
##### A fictitious Kubernetes Operator

This is an exercise I took up to understand the Kubernetes Operator framework better. The fictitious app comprises of a cluster of nginx servers.The following rules always apply for the clusters:
1. Servers are either masters or slaves or neither.
2. Slaves belong to two categories - "category 1 slaves" and "category 2 slaves"
3. The cluster size is configurable through the Ikshvaku CRD
4. If the cluster size is 1 or 2, there are no masters nor slaves.
5. If the cluster size is even, there are two masters and the rest of the servers are equally distributed among "category 1" and "category 2"
6. If the cluster size is odd, there's one master server and the rest of the servers are equally distributed among "category 1" abd "category 2"

