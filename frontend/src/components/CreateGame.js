import { useState } from "react"
import Popup from "./Popup"
import { useNavigate } from 'react-router-dom';
import './CreateGame.css'


const getUrl = async () => {
    const response = await fetch(`${window.location.origin}/api/service`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ action: 'createGame' }),
    });

    if (!response.ok) {
        throw new Error('Network response was not ok');
    }

    const data = await response.json();
    return "/" + data.gameId;
};


const CreateGame = () => {
    const navigate = useNavigate();
    const [seconds, setSeconds] = useState(null);
    const [url, setUrl] = useState(null);

    const createGame = async () => {
        const url = await getUrl()
        setUrl(url)
        setSeconds(5);

        const currentDomain = window.location.origin;
        await navigator.clipboard.writeText(currentDomain+url);

        setTimeout(() => {
            setTimeout( () => {navigate(url)}, 5000);
        }, 500)
        
    }

    return <div className="bg-image">
        <div className="blurred-background p-5 mb-4 rounded-3">
            <div className="container-fluid py-5">
                <h1 className="display-5 fw-bold">Welcome to planing poker!</h1>
                <p className="col-md-8 fs-4">To start press button:</p>
                <button onClick={createGame} id="createLinkButton" className="btn btn-primary btn-lg" type="button">
                    Create new game <i className="bi bi-play-fill"></i>
                </button>
            </div>
        </div>

    { seconds ? <Popup header={"New game created!"}><p><h1><i className="bi bi-clipboard-check"></i></h1>{"Link copied to clipboard, you will be redirected in " + seconds + " seconds or click "} <a href={url}>here</a> </p></Popup> : <></> }

    </div>

}

export default CreateGame