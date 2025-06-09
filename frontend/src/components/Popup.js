import "./Popup.css"

const Popup = ({header, text}) => {
    return <div className="popup-wrapper pt-5">
    <div className="popup-content container-sm">
      <h4>{header}</h4>
      <p>{text}</p>
    </div>
  </div>
}

export default Popup