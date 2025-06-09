import { useState } from 'react';

const GetUserName = ({ setUserName }) => {
  const [userName, setUserNameLocal] = useState('');

  const handleSubmit = (event) => {
    event.preventDefault();
    setUserName(userName);
  };

  return (
    <div className="row justify-content-center">
      <div className="col-md-6">
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label htmlFor="userName" className="form-label">
              Enter your username
            </label>
            <input
              type="text"
              className="form-control"
              id="userName"
              value={userName}
              onChange={(e) => setUserNameLocal(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn btn-primary">
            Start
          </button>
        </form>
      </div>
    </div>
  );
};

export default GetUserName;
