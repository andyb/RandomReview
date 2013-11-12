Random Review deployment
=================

Deployed to Rackspace Cloud - SSLNOPSREVIEW001 instance
- ssh into instance.
- binary into home dir.
- to deploy use SCP to transfer files e.g scp random_review root@162.13.146.76 and scp reviewers.json root@162.13.146.76
- Ubuntu upstart script used to manage lifecycle. This is located in /etc/init/review.conf
	- To start: sudo start review
	- To stop: sudo stop review

- logs for app in /var/log/upstart/review.log