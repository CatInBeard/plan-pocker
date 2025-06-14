import "./Popup.css"

const Popup = ({header, children, closeAction = null}) => {
    return <div className="popup-wrapper pt-5">
    <div className="popup-content container-sm">
      {header && <><h4 className="d-flex justify-content-between align-items-center"><span>{header}</span> {closeAction && <i onClick={closeAction} className="bi bi-x-lg float-right"></i>}</h4>
      <hr/></>}
      {children}
    </div>
  </div>
}

export default Popup