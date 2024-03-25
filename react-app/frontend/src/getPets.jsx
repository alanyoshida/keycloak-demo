import { useState, useEffect } from 'react';

const GetPets = () => {
    const [pets, setPets] = useState([]);
    useEffect(() => {
        var front_access_token = localStorage.getItem("token");

        var options = {
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + front_access_token,
            }
        }
        fetch('http://localhost:4000/protected/pets/list', options)
        .then(response => {
            return response.json()
        })
        .then((pets) => {
            console.log("SUCESSO", pets);
            setPets(pets);
        })
        .catch(error => console.error(error));

    }, []);
    return (
        <table border="1">
        <tbody>
        <tr>
            <th>Nome</th><th>Raca</th><th>Tipo</th><th>Idade</th><th>genero</th>
        </tr>
        {/* { pets ? JSON.stringify(pets, null, 2): "" } */}
        { pets && pets.pets ?
             pets.pets.map((el) => {
                 return <tr>
                  <td>{el.nome}</td>
                  <td>{el.raca}</td>
                  <td>{el.tipo}</td>
                  <td>{el.idade}</td>
                  <td>{el.genero}</td>
                  </tr>
              }): ""
         }
        </tbody>
      </table>
    );
};
export default GetPets