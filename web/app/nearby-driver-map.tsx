"use client";
import {API_URL, DEFAULT_LATITUDE, DEFAULT_LONGITUDE} from "@/lib/constants";
import {
  Driver,
  Location as LocationType,
  RequestRideProps,
  RouteInfo,
} from "@/types";
import L from "leaflet";
import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";
import {MapContainer, Marker, Popup, Rectangle, TileLayer} from "react-leaflet";
import {useNearbyDrivers} from "@/hooks/use-nearby-drivers";
import {useRef, useState} from "react";
import {decodeGeoHash} from "@/lib/geo-hash";
import {Button} from "@/components/ui/button";
import {DriverList} from "./driver-list";
import {RoutingControl} from "./routing-control";
import {MapClickHandler} from "./map-click-handler";

// Default icon for the map

const DefaultIcon = L.icon({
  iconUrl: icon.src,
  shadowUrl: iconShadow.src,
  iconSize: [25, 41], // size of the icon
  iconAnchor: [12, 41], // where the icon will appear
  popupAnchor: [1, -34], // where the popup will appear
  shadowSize: [41, 41], // size of the shadow
});

// Set the default icon for the map
L.marker.prototype.options.icon = DefaultIcon;

export default function NearbyDriversMap() {
  const [selectedDriver, setSelectedDriver] = useState<Driver | null>(null);
  const [route, setRoute] = useState<[number, number][]>([]);
  const mapRef = useRef<L.Map>(null);

  const location: LocationType = {
    latitude: DEFAULT_LATITUDE,
    longitude: DEFAULT_LONGITUDE,
  };

  const [destination, setDestination] = useState<[number, number] | null>(null);
  const {drivers, error} = useNearbyDrivers(location);

  const userMarker = new L.Icon({
    iconUrl:
      "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ed/Map_pin_icon.svg/176px-Map_pin_icon.svg.png",
    iconSize: [40, 40], // Size of the marker
    iconAnchor: [20, 40], // Anchor point
  });

  const driverMarker = new L.Icon({
    iconUrl: "https://www.svgrepo.com/show/25407/car.svg",
    iconSize: [30, 30],
    iconAnchor: [15, 30],
  });

  // function to get the bounds of a geohash
  const getGeohashBounds = (geohash: string): [number, number][] => {
    const {
      latitude: [minLat, maxLat],
      longitude: [minLng, maxLng],
    } = decodeGeoHash(geohash);
    return [
      [minLat, minLng] as [number, number],
      [maxLat, maxLng] as [number, number],
    ];
  };

  const handleMapClick = async (e: L.LeafletMouseEvent) => {
    setDestination([e.latlng.lat, e.latlng.lng]);

    await requestRide({
      pickup: [location.latitude, location.longitude],
      destination: [e.latlng.lat, e.latlng.lng],
    });
  };

  const requestRide = async (props: RequestRideProps): Promise<RouteInfo> => {
    const {pickup, destination} = props;
    const payload = {
      pickup: {
        latitude: pickup[0],
        longitude: pickup[1],
      },
      destination: {
        latitude: destination[0],
        longitude: destination[1],
      },
    };

    const response = await fetch(`${API_URL}/trip`, {
      method: "POST",
      body: JSON.stringify(payload),
    });
    const data = (await response.json()) as RouteInfo;
    const route = data.route;
    const parsedRoute = route.geometry[0].coordinates.map(
      (coord) => [coord.longitude, coord.latitude] as [number, number]
    );

    setRoute(parsedRoute);
    return data;
  };

  if (error) {
    return <div>Error: {error}</div>;
  }
  return (
    <MapContainer
      center={[location.latitude, location.longitude]}
      zoom={13}
      style={{
        height: "100vh",
        width: "100%",
      }}
      ref={mapRef}
    >
      <TileLayer
        url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
      />
      <Marker
        position={[location.latitude, location.longitude]}
        icon={userMarker}
      >
        {/* render geo-hash grid cells */}
        {drivers.map((driver) => (
          <Rectangle
            key={`grid-${driver.geohash}`}
            bounds={getGeohashBounds(driver.geohash)}
            pathOptions={{
              color: "#3388ff",
              weight: 1,
              fillOpacity: 0.1,
              fillColor: "#3388ff",
            }}
          >
            <Popup>Geohash: {driver.geohash}</Popup>
          </Rectangle>
        ))}

        {/* render driver markers */}

        {drivers.map((driver) => (
          <Marker
            key={driver.driver_id}
            position={[driver.location.latitude, driver.location.longitude]}
            icon={driverMarker}
          >
            <Popup>
              Driver ID: {driver.driver_id}
              <br />
              Geohash: {driver.geohash}
            </Popup>
          </Marker>
        ))}
        {destination && (
          <Marker position={destination} icon={userMarker}>
            <Popup>Destination</Popup>
          </Marker>
        )}

        {destination ? (
          <DriverList drivers={drivers} onSelectDriver={setSelectedDriver} />
        ) : (
          <div className="flex items-center justify-center h-full">
            <p className="text-lg font-semibold text-gray-500">
              Click on the map to set a destination
            </p>
          </div>
        )}
        {selectedDriver && (
          <div className="mt-4 z-[9999] absolute bottom-0 right-0">
            <Button className="w-full">
              Request Ride with {selectedDriver.driver_id}
            </Button>
          </div>
        )}

        <RoutingControl route={route} />
        <MapClickHandler onClick={handleMapClick} />
      </Marker>
    </MapContainer>
  );
}
