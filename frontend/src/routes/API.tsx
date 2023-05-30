// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// import ReactJson from 'react-json-view'

function APIState() {
	return <Container>
		<ReactJson src={lastJsonMessage} />
	</Container>
}