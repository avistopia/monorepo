.PHONY: init-ansible
init-ansible:
	ansible-galaxy collection install community.docker --force
	ansible-galaxy install geerlingguy.docker

.PHONY: ansible-setup-instance
ansible-setup-instance:
	ansible-playbook -i ansible/inventory.ini ansible/setup_instance.yaml

.PHONY: ansible-deploy
ansible-deploy:
	ansible-playbook -i ansible/inventory.ini ansible/deploy/$(app).yaml -e "env_file_path=../../../$(app)/production.env"
