
const OKAY = 'okay'
const REJECT = 'reject'

// HELO annd FROM are both strings.
// RCPT is an array of strings.
// Both FROM and RCPT have been set to lowercase already.

function execute({ HELO, FROM, RCPT }) {
	if (HELO == 'localhost') {
		return REJECT
	}
	if (!['joe@example.com', 'mary@example.com'].includes(FROM)) {
		return REJECT
	}
	if (RCPT.includes('joe@example.com')) {
		return REJECT
	}
	return OKAY
}
