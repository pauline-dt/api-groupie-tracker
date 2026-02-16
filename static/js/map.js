
let map;
let markers = [];
let geocoder;


function initMap() {
   
    const defaultCenter = { lat: 48.8566, lng: 2.3522 };

    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 2,
        center: defaultCenter,
        styles: [
            {
                "elementType": "geometry",
                "stylers": [{ "color": "#1e293b" }]
            },
            {
                "elementType": "labels.text.stroke",
                "stylers": [{ "color": "#0f172a" }]
            },
            {
                "elementType": "labels.text.fill",
                "stylers": [{ "color": "#94a3b8" }]
            },
            {
                "featureType": "water",
                "elementType": "geometry",
                "stylers": [{ "color": "#334155" }]
            }
        ]
    });

    geocoder = new google.maps.Geocoder();

    // Géocoder et afficher tous les lieux de concert
    if (typeof artistData !== 'undefined' && artistData.locations) {
        geocodeLocations(artistData.locations);
    }
}


async function geocodeLocations(locations) {
    const bounds = new google.maps.LatLngBounds();

    for (const locationData of locations) {
        try {
            const position = await geocodeAddress(locationData.location);
            
            if (position) {
                addMarker(position, locationData);
                bounds.extend(position);
            }
        } catch (error) {
            console.error(`Erreur de géocodage pour ${locationData.location}:`, error);
        }
    }


    if (markers.length > 0) {
        map.fitBounds(bounds);
    }
}


function geocodeAddress(address) {
    return new Promise((resolve, reject) => {
        // Nettoyer l'adresse
        const cleanAddress = address.replace(/_/g, ' ').replace(/-/g, ' ');

        geocoder.geocode({ address: cleanAddress }, (results, status) => {
            if (status === 'OK' && results[0]) {
                resolve(results[0].geometry.location);
            } else {
                reject(new Error(`Geocoding failed: ${status}`));
            }
        });
    });
}

// Ajouter un marqueur sur la carte
function addMarker(position, locationData) {
    const marker = new google.maps.Marker({
        position: position,
        map: map,
        title: locationData.location.replace(/_/g, ' ').replace(/-/g, ' '),
        animation: google.maps.Animation.DROP
    });


    const dates = locationData.dates.map(date => `<li>${date}</li>`).join('');
    const contentString = `
        <div style="color: #0f172a; padding: 10px; max-width: 300px;">
            <h3 style="margin-top: 0; color: #6366f1;">📍 ${marker.title}</h3>
            <h4 style="margin-top: 10px; margin-bottom: 5px;">Dates des concerts:</h4>
            <ul style="margin: 5px 0; padding-left: 20px;">
                ${dates}
            </ul>
        </div>
    `;

    const infowindow = new google.maps.InfoWindow({
        content: contentString
    });


    marker.addListener('click', () => {

        markers.forEach(m => {
            if (m.infowindow) {
                m.infowindow.close();
            }
        });
        infowindow.open(map, marker);
    });

    marker.infowindow = infowindow;
    markers.push(marker);
}


if (!document.getElementById('map')) {
    console.log('Carte désactivée - élément non trouvé');
} else if (typeof google === 'undefined') {
    document.getElementById('map').innerHTML = `
        <div style="display: flex; align-items: center; justify-content: center; height: 100%; background: var(--background); color: var(--text-secondary); text-align: center; padding: 2rem;">
            <div>
                <p style="font-size: 1.2rem; margin-bottom: 1rem;">🗺️ Carte non disponible</p>
                <p style="font-size: 0.9rem;">Pour activer la géolocalisation, vous devez configurer une clé API Google Maps.</p>
                <p style="font-size: 0.8rem; margin-top: 1rem;">Les lieux de concerts sont listés ci-dessous.</p>
            </div>
        </div>
    `;
}


function clearMarkers() {
    markers.forEach(marker => marker.setMap(null));
    markers = [];
}