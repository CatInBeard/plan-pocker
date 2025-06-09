import GetUserName from "./GetUserName";
import Header from "./Header";
import Container from "./Container";
import PlayerTable from "./PlayerTable";
import WsClient from "../Websocket/WsClient";
import CardPicker from "./CardPicker";
import Settings from "./Settings";
import style from "./Game.module.css"

import { useState, useEffect, useRef } from 'react';


const Game = ({ id }) => {

    const [userName, setUserName] = useState(null);
    const [voters, setVoters] = useState([])
    const [deck, setDeck] = useState([])
    const [allowCustomDeck, setAllowCustomDeck] = useState(false)
    const [revealedValue, setRevealedValue] = useState(null)
    const [showSettings, setShowSettings] = useState(false)

    const openSettings = () => {
        setShowSettings(true)
    }

    const closeSettings = () => {
        setShowSettings(false)
    }

    let [wsClient, setWsClient] = useState(null)
    let [connectionStatus, setConncetionStatus] = useState("Not set")

    let [selectedCard, setSelectedCard] = useState(null);

    let onWsMessage = (message) => {

        switch(message.action){
            case "voters":
                setVoters(message.voters)
            break;
            case "deck":
                setDeck(message.deck)
                setAllowCustomDeck(message.allowCustom)
            break;
        }
    }

    useEffect(() => {
        if (!wsClient) {
            setWsClient(new WsClient(setConncetionStatus, onWsMessage))
        }
    }, [])

    const updateValueFromLocalStorage = () => {
        const storedValue = localStorage.getItem('userName');
        if (storedValue !== null) {
            setUserName(storedValue);
        }
    };

    useEffect(() => {
        updateValueFromLocalStorage();
    }, []);

    useEffect(() => {
        if(userName !== null && String(userName).length > 0){
            localStorage.setItem('userName', userName);
        }
    }, [userName]);

    useEffect(() => {
        const handleStorageChange = (event) => {
            if (event.key === 'userName') {
                updateValueFromLocalStorage();
            }
        };

        window.addEventListener('storage', handleStorageChange);

        return () => {
            window.removeEventListener('storage', handleStorageChange);
        };
    })

     useEffect(() => {
        let intervalId;

        if (connectionStatus === 'established' && userName) {
        const sendMessage = () => {
            const message = {
            userName: userName,
            action: 'connected',
            };
            wsClient.send(JSON.stringify(message))
        };

        sendMessage();

        intervalId = setInterval(sendMessage, 2000);
        }

        return () => {
        if (intervalId) {
            clearInterval(intervalId);
        }
        };
    }, [userName, connectionStatus, wsClient]);


    if(!userName){
        return <>
            <Header />
            <Container>
                <GetUserName setUserName={setUserName}/>
            </Container>
        </>
    }

    const revealAction = () => {
        wsClient.send(
            JSON.stringify({"action": "reveal"})
        );
    }

    const pickCB = (pick) => {
        setSelectedCard(pick)
    }

    return <>
        <Header>Planing | connection status: {connectionStatus}
            <div className={style.settings} onClick={openSettings}>
                <i className="bi bi-gear"></i>
            </div>
        </Header>
        <Container>
            {showSettings && <Settings/>}
            <PlayerTable voters={voters} revealedValue={revealedValue} revealAction={revealAction}/>
            <CardPicker availableCards={deck} allowCustom={allowCustomDeck} selectCallback={pickCB} selectedCard={selectedCard} />
        </Container>
    </>
}

export default Game