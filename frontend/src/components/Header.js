function Header({children = "Planing"}) {

  return (
    <header className="container border-bottom pb-4 pt-1">
      {children}
    </header>
  );
}

export default Header;
