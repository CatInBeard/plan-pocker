import { useState } from "react"
import Popup from "./Popup"
import { useNavigate } from 'react-router-dom';

function generateRandomString(len) {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    const randomValues = new Uint32Array(10);
    window.crypto.getRandomValues(randomValues);

    for (let i = 0; i < len; i++) {
        result += chars[randomValues[i] % chars.length];
    }

    return result;
}


const getUrl = () => {
    return "/" +  generateRandomString(5)
}

const CreateGame = () => {
    const navigate = useNavigate();
    const [seconds, setSeconds] = useState(null);
    const [url, setUrl] = useState(null);

    const createGame = async () => {
        const url = getUrl()
        setUrl(url)
        setSeconds(5);

        const currentDomain = window.location.origin;
        await navigator.clipboard.writeText(currentDomain+url);

        setTimeout(() => {
            setTimeout( () => {navigate(url)}, 5000);
        }, 500)
        
    }

    return <div>
        <div className="p-5 mb-4 bg-light rounded-3">
            <div className="container-fluid py-5">
                <h1 className="display-5 fw-bold">Welcome to planing poker!</h1>
                <p className="col-md-8 fs-4">You can create new game:</p>
                <button onClick={createGame} id="createLinkButton" className="btn btn-primary btn-lg" type="button">
                    Create a game <i className="bi bi-play-fill"></i>
                </button>
            </div>
        </div>

    { seconds ? <Popup header={"New game created!"} text={<p><h1><i className="bi bi-clipboard-check"></i></h1>{"Link copied to clipboard, you will be redirected in " + seconds + " seconds or click "} <a href={url}>here</a> </p>}/> : <></> }

    </div>

}

export default CreateGame