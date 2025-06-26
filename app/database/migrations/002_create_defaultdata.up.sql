INSERT OR IGNORE INTO types (id, singular, plural, icon, schema) VALUES ('alert', 'Alert', 'Alerts', 'AlertTriangle', '{"type": "object", "properties": { "severity": { "title": "Severity", "enum":  ["Low", "Medium", "High"]}}, "required": ["severity"]}');
INSERT OR IGNORE INTO types (id, singular, plural, icon, schema) VALUES ('incident', 'Incident', 'Incidents', 'Flame', '{"type": "object", "properties": { "severity": { "title": "Severity", "enum":  ["Low", "Medium", "High"]}}, "required": ["severity"]}');
INSERT OR IGNORE INTO types (id, singular, plural, icon, schema) VALUES ('vulnerability', 'Vulnerability', 'Vulnerabilities', 'Bug', '{"type": "object", "properties": { "severity": { "title": "Severity", "enum":  ["Low", "Medium", "High"]}}, "required": ["severity"]}');

INSERT OR IGNORE INTO users (id, name, username, passwordHash, tokenKey, active, admin) VALUES ('system', 'System', 'system', '', lower(hex(randomblob(26))), true, true);
