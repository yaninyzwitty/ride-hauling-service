export interface Location {
    latitude: number;
    longitude: number;
}

export interface Driver {
    driver_id: string;
    location: Location;
    geohash: string;
}

export type RequestRideProps =  {
    pickup: [number, number],
    destination: [number, number],
}

export interface RouteInfo {
    route: {
        geometry: {
            coordinates: Location[]
        }[]
    }
}