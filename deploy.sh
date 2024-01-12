ansible-playbook -i playbooks/hosts.ini playbooks/ais_datafs_mkfs.yml -e ais_hosts=cpu-worker-node  --become

# only for ubuntu < 18.04
# ansible-playbook -i playbooks/hosts.ini playbooks/ais_enable_multiqueue.yml -e playhosts=cpu-worker-node --become

ansible-playbook -i playbooks/hosts.ini playbooks/ais_deploy_operator.yml
ansible-playbook -i playbooks/hosts.ini playbooks/ais_deploy_cluster.yml -e cluster="ais-1"
# ansible-playbook -i playbooks/hosts.ini playbooks/ais_destroy_cluster.yml -e cluster="ais-1"

ansible-playbook -i playbooks/hosts.ini playbooks/ais_cleanup_all.yml -e cluster="ais-1"