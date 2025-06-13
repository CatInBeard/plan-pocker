import GetUserName from "./GetUserName";
import Header from "./Header";
import Container from "./Container";
import PlayerTable from "./PlayerTable";
import WsClient from "../Websocket/WsClient";
import CardPicker from "./CardPicker";
import Settings from "./Settings";
import style from "./Game.module.css"
import { v4 as uuidv4 } from 'uuid';

import { useState, useEffect, useRef } from 'react';


const Game = ({ id }) => {

    const [userName, setUserName] = useState(null);
    const [voters, setVoters] = useState([])
    const [deck, setDeck] = useState([])
    const [allowCustomDeck, setAllowCustomDeck] = useState(false)
    const [revealedValue, setRevealedValue] = useState(null)
    const [showSettings, setShowSettings] = useState(false)
    const [uid, setUid] = useState(null);
    const [selectedCard, setSelectedCard] = useState(null);
    const [wsClient, setWsClient] = useState(null)
    const [connectionStatus, setConnectionStatus] = useState("Not set")

    const openSettings = () => {
        setShowSettings(true)
    }

    const setSettings = (deck, isCustomDeckAllowed) => {
        const message = {
            gameId: id,
            action : 'setSettings',
            deck : deck.map(item => Number(item)),
            uid : uid,
            isCustomDeckAllowed: isCustomDeckAllowed
        };
        wsClient.send(message)
        setShowSettings(false)
    }

    const selectedCardRef = useRef(selectedCard);
    useEffect(() => {
        selectedCardRef.current = selectedCard;
    }, [selectedCard]);

    let onWsMessage = (message) => {

        switch(message.action){
            case "voters":
                setVoters(message.voters ?? [])
                setRevealedValue(message.vote)
                if(message.vote){
                    setSelectedCard(null)
                }
            break;
            case "deck":
                setDeck(message.deck)
                setAllowCustomDeck(message.allowCustom)
            break;
        }
    }

    useEffect(() => {
        if (!wsClient) {
            setWsClient(new WsClient(setConnectionStatus, onWsMessage))
        }
    }, [])

    const updateValueFromLocalStorage = () => {
        const storedValue = localStorage.getItem('userName');
        if (storedValue !== null) {
            setUserName(storedValue);
        }
        let uid = localStorage.getItem('uid');
        if (uid !== null) {
            setUid(uid);
        } else {
            uid = uuidv4();
            setUid(uid);
            localStorage.setItem('uid', uid);
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
            wsClient.send({
            userName: userName,
            uid: uid,
            gameId: id,
            action: 'connect',
            vote: selectedCardRef.current ?? 0
            })
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
                {
                    action: "reveal",
                    gameId: id,
                    uid: uid
                }
        );
    }

    const startNewAction = () => {
        wsClient.send(
                {
                    action: "start",
                    gameId: id,
                    uid: uid
                }
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
            {showSettings && <Settings deck={deck} allowCustomDeck={allowCustomDeck} setSettings={setSettings} />}
            <PlayerTable voters={voters} revealedValue={revealedValue} revealAction={revealAction} startNewAction={startNewAction} />
            <CardPicker availableCards={deck} allowCustom={allowCustomDeck} selectCallback={pickCB} selectedCard={selectedCard} />
        </Container>
    </>
}

export default Game