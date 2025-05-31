package fakedata

import "github.com/brianvoe/gofakeit/v7"

func fakeTicketDescription() string {
	return random([]string{
		"Unauthorized access attempt detected in the main server room.",
		"Multiple failed login attempts from an unknown IP address.",
		"Suspicious file download flagged by antivirus software.",
		"User account locked due to repeated incorrect password entries.",
		"Unusual network activity observed on the internal firewall.",
		"Phishing email reported by several employees.",
		"Sensitive data transfer detected outside the approved hours.",
		"Malware infection found on a workstation in the finance department.",
		"Unauthorized device connected to the company network.",
		"Brute-force attack attempt on the admin account detected.",
		"Security patch required for vulnerability in outdated software.",
		"External IP address attempting to probe network ports.",
		"Suspicious behavior detected by user in HR department.",
		"Unauthorized software installation on company laptop.",
		"Access control system malfunction at the main entrance.",
		"DDoS attack detected on company web server.",
		"Unusual outbound traffic to a known malicious domain.",
		"Potential insider threat flagged by behavior analysis tool.",
		"Compromised credentials detected on dark web.",
		"Encryption key rotation required for compliance with security policy.",
	})
}

func fakeTicketComment() string {
	return random([]string{
		"Ticket opened by user.",
		"Initial investigation started.",
		"Further analysis required.",
		"Escalated to security team.",
		"Action taken to mitigate risk.",
		"Resolution in progress.",
		"User notified of incident.",
		"Security incident confirmed.",
		"Containment measures implemented.",
		"Root cause analysis underway.",
		"Forensic investigation initiated.",
		"Data breach confirmed.",
		"Incident response team activated.",
		"Legal counsel consulted.",
		"Public relations notified.",
		"Regulatory authorities informed.",
		"Compensation plan developed.",
		"Press release drafted.",
		"Media monitoring in progress.",
		"Post-incident review scheduled.",
	})
}

func fakeTicketTimelineMessage() string {
	return random([]string{
		"Initial investigation started.",
		"Further analysis required.",
		"Escalated to security team.",
		"Action taken to mitigate risk.",
		"Resolution in progress.",
		"User notified of incident.",
		"Security incident confirmed.",
		"Containment measures implemented.",
		"Root cause analysis underway.",
		"Forensic investigation initiated.",
		"Data breach confirmed.",
		"Incident response team activated.",
		"Legal counsel consulted.",
		"Public relations notified.",
		"Regulatory authorities informed.",
		"Compensation plan developed.",
		"Press release drafted.",
		"Media monitoring in progress.",
		"Post-incident review scheduled.",
	})
}

func fakeTicketTask() string {
	return random([]string{
		"Interview witnesses.",
		"Review security camera footage.",
		"Analyze network traffic logs.",
		"Scan for malware on affected systems.",
		"Check for unauthorized software installations.",
		"Conduct vulnerability assessment.",
		"Implement security patch.",
		"Change firewall rules.",
		"Reset compromised credentials.",
		"Isolate infected systems.",
		"Monitor for further suspicious activity.",
		"Coordinate with law enforcement.",
		"Notify affected customers.",
		"Prepare incident report.",
		"Update security policies.",
		"Train employees on security best practices.",
		"Conduct post-incident review.",
		"Implement lessons learned.",
		"Improve incident response procedures.",
		"Enhance security awareness program.",
	})
}

func random[T any](e []T) T {
	return e[gofakeit.IntN(len(e))]
}
