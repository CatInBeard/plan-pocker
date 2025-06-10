import "./Popup.css"

const Popup = ({header, children}) => {
    return <div className="popup-wrapper pt-5">
    <div className="popup-content container-sm">
      {header && <><h4>{header}</h4>
      <hr/></>}
      {children}
    </div>
  </div>
}

export default Popup