import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

const Login = (props) => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [usernameError, setUsernameError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  const navigate = useNavigate()

  const onButtonClick = () => {
    // Set initial error values to empty
    setUsernameError('')
    setPasswordError('')

    // Check if the user has entered both fields correctly
    if ('' === username) {
        setUsernameError('Please enter your user')
      return
    }

    if ('' === password) {
      setPasswordError('Please enter a password')
      return
    }

    if (password.length < 7) {
      setPasswordError('The password must be 8 characters or longer')
      return
    }

    logIn()

    // Authentication calls will be made here...
  }

  const logIn = () => {
    const formData = new URLSearchParams();
    formData.append('grant_type', 'password');
    formData.append('client_id', 'frontend');
    formData.append('username', username);
    formData.append('password', password);

    fetch('http://keycloak.default.svc.cluster.local/realms/master/protocol/openid-connect/token', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: formData.toString(),
    })
    .then(response => {
        if (200 === response.status) {
            localStorage.setItem('user', JSON.stringify({ username }));
            props.setLoggedIn(true);
            props.setUsername(username);
        } else {
            window.alert('Wrong user or password');
        }
        return response.json()
    })
    .then((r) => {
        console.log(r)
        console.log(r.access_token)
        localStorage.setItem('token', r.access_token);
        navigate('/');
    })
  }
  return (
    <div className={'mainContainer'}>
      <div className={'titleContainer'}>
        <div>Login</div>
      </div>
      <br />
      <div className={'inputContainer'}>
        <input
          value={username}
          placeholder="Enter your user here"
          onChange={(ev) => setUsername(ev.target.value)}
          className={'inputBox'}
        />
        <label className="errorLabel">{usernameError}</label>
      </div>
      <br />
      <div className={'inputContainer'}>
        <input
          type='password'
          value={password}
          placeholder="Enter your password here"
          onChange={(ev) => setPassword(ev.target.value)}
          className={'inputBox'}
        />
        <label className="errorLabel">{passwordError}</label>
      </div>
      <br />
      <div className={'inputContainer'}>
        <input className={'inputButton'} type="button" onClick={onButtonClick} value={'Log in'} />
      </div>
    </div>
  )
}

export default Login