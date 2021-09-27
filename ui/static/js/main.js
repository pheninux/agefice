var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

function openModal(target) {

    switch (target) {
        case "formation" :
            document.getElementById('modal-formation').style.display='block' ;
            break
        case "entreprise" :
            document.getElementById('modal-entreprise').style.display='block' ;
            break
        case "document" :
            document.getElementById('modal-document').style.display='block' ;
            break
    }
}



function nouvelleFormation() {

	document.getElementById("dateDeb").value = "" ;
	document.getElementById("dateFin").value = "" ;
	document.getElementById("nbrHeures").value = "" ;
	document.getElementById("cout").value = "" ;
	document.getElementById('intitule').style.visibility = 'visible';

}
