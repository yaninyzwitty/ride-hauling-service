import Geohash from 'latlon-geohash';

export function decodeGeoHash(geohash: string) {
  const bounds = Geohash.bounds(geohash);
  return {
    latitude: [bounds.sw.lat, bounds.ne.lat],
    longitude: [bounds.sw.lon, bounds.ne.lon],
  };
}