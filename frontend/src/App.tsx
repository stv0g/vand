// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

import { Navbar, Nav, NavDropdown, Container } from 'react-bootstrap';
// import useWebSocket from 'react-use-websocket';

import logo from './logo.svg';
import './App.scss';

function App() {
	// const socketUrl = '/api/v1/ws';
	// const {
	// 	sendMessage,
	// 	sendJsonMessage,
	// 	lastMessage,
	// 	lastJsonMessage,
	// 	readyState,
	// 	getWebSocket,
	// } = useWebSocket(socketUrl, {
	// 	onOpen: () => console.log('opened'),
	// 	// Will attempt to reconnect on all close events, such as server shutting down
	// 	shouldReconnect: (closeEvent) => true,
	// });

	return (
		<div className="App">
			<Navbar bg="light" expand="lg">
				<Container>
					<Navbar.Brand href="#home">
						<img alt="VANd logo" src={logo} height="30" className="d-inline-block align-top" /> VANd
					</Navbar.Brand>
					<Navbar.Toggle aria-controls="basic-navbar-nav" />
					<Navbar.Collapse id="basic-navbar-nav">
						<Nav className="me-auto">
							<Nav.Link href="#home">Home</Nav.Link>
							<Nav.Link href="https://tracks.0l.de/?user=bus&layers=last,line">Tracks</Nav.Link>
							<NavDropdown title="About" id="basic-nav-about">
								<NavDropdown.Item href="https://github.com/stv0g/vand">GitHub</NavDropdown.Item>
								<NavDropdown.Item href="https://www.steffenvogel.de">@stv0g</NavDropdown.Item>
							</NavDropdown>
						</Nav>
					</Navbar.Collapse>
				</Container>
			</Navbar>
			<Container>
			</Container>
		</div>
	);
}

export default App;
