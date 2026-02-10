# React

React est une bibliothèque JavaScript frontale à code source ouvert permettant de créer
des interfaces utilisateur ou des composants d'interface utilisateur.

# Formation

https://www.youtube.com/watch?v=f0X1Tl8aHtA&t=36s
https://react.dev/learn/build-a-react-app-from-scratch
https://www.youtube.com/watch?v=WQhQGWjreVg&list=PL4krRxGtKDzkiOZ7Ypfg-yGFa7T6Ahr5m

    Props une façon d'envoyer des paramètres à des composants
    hook : fonction proposée par react --> tous les hook commence par use.. exemple useState()
    	state : permet de suivre une variable, et a chaque qu'elle est modifiée, il va réactualisé le composant !

    destructuring
    	En React, on l'utilise tout le temps pour récupérer les props ou les données d'un état.
    	Sans destructuring :
    		On doit répéter props à chaque fois.
    			function Profil(props) {
    			  return <h1>Bonjour, {props.nom} ! Tu as {props.age} ans.</h1>;
    			}
    	Avec destructuring :
    		On "extrait" directement les propriétés de l'objet props.
    			function Profil({ nom, age }) {
    			  return <h1>Bonjour, {nom} ! Tu as {age} ans.</h1>;
    			}
    Le destructuring de Tableaux
    	C'est la base de l'utilisation des Hooks comme useState.
    	Lorsqu'une fonction renvoie un tableau, le destructuring permet de nommer les éléments selon leur position.
    	// useState renvoie un tableau : [valeur, fonction_de_mise_a_jour]

    		Exemple
    		function donnerCouleurs() {
    		  return ["bleu", "rouge"];
    		}

    		// Sans destructuring :
    		const resultat = donnerCouleurs();
    		const premier = resultat[0]; // "bleu"
    		const second = resultat[1];  // "rouge"

    		// AVEC destructuring :
    		const [premier, second] = donnerCouleurs();
    		// premier devient "bleu"
    		// second devient "rouge"
    		Le signe [ ] à gauche du = dit à JavaScript :
    			"Prends le premier élément du tableau et mets-le dans la variable 'premier', puis prends le deuxième et mets-le dans 'second'".

    	La fonction useState de React fonctionne exactement de cette manière.
    		Elle renvoie toujours un tableau contenant exactement deux éléments :
    		* La valeur actuelle de l'état.
    		* Une fonction pour modifier cette valeur.

    			const etatTableau = useState(0); // Renvoie par ex: [0, function]
    				const count = etatTableau[0];    // On récupère la valeur
    				const setCount = etatTableau[1]; // On récupère la fonction


    			En écrivant const [count, setCount] = useState(0);, tu fais ces deux étapes en une seule ligne.
    				const [pomme, changerPomme] = useState(0);


    Apps.jsx
    	let number

    	const [number, setNumber] = useState(0)
    	function add() {
    		// number++ :: plus le droit à cause de useState !!
    		setNumber(number+1)
    		console.log(number) <-- Ancienne valeur !!
    	}
    	function substract {
    		// number--  :: Plus le droit à cause de useState !!
    		setNumber(number+1)
    		console.log(number)
    	}

    	return {
    		<div>
    			...
    			<Button content = "-" click={substract} />
    			<Button content = "+" click={add} />
    			<Counter> {nombre} </Counter> --> Autre façon de passer un attribu
    			...
    		</div>
    	}

    		import { useState, useEffect } from 'react';

    		function MonComposant() {
    		  const [number, setNumber] = useState(0);

    		  // Ce code s'exécute APRÈS que le composant a été mis à jour
    		  useEffect(() => {
    			console.log("La nouvelle valeur est :", number);
    		  }, [number]); // <--- On dit à React de surveiller "number"

    		  function add() {
    			setNumber(number + 1);
    			// Ici, le console.log afficherait encore l'ancienne valeur
    		  }

    		  return (
    			<button onClick={add}>Ajouter +1</button>
    		  );
    		}

    Button.jsx
    	export default function Button(props) {
    		console.log(props.content)
    		//alert("Montre une fenetre")
    		const nbut = props.content
    		return <button> {nbut} </button>
    	}

    	ou plus simplement

    	export default function Button( {content, click}) {
    		return <button onClick={click}> {content} </button>
    	}

    Counter.jsx
    	export default function Counter( {childre}) {
