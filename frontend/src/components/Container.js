import "./Container.css"

const Container = ({children}) => {

    return <main className="container border-left border-right container-height">
        {children}
    </main>

}

export default Container