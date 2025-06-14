import { WEBSOCKET_URL } from '@/lib/constants';
import { Driver, Location as LocationType } from '@/types';
import { useState, useEffect } from 'react';




/**
 * Hook to manage real-time nearby drivers using WebSocket
 * @param location - Current user location
 * @returns Object containing drivers array and error state
 */


export const useNearbyDrivers = (location: LocationType) => {
    const [drivers, setDrivers] = useState<Driver[]>([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        // Skip if no location provided
        if (!location) return;

        // Initialize WebSocket connection
        const ws = new WebSocket(`${WEBSOCKET_URL}/drivers`);
        const stringifiedLocation = JSON.stringify(location);

        // Connection opened
        ws.onopen = () => {
            ws.send(stringifiedLocation);
        };

        // Handle incoming driver updates
        ws.onmessage = (event) => {
            try {
                const drivers = JSON.parse(event.data) as Driver[];
                setDrivers(drivers);
            } catch (err) {
                setError('Failed to parse driver data');
                console.error('Parse error:', err);
            }
        };

        // Handle connection closure
        ws.onclose = () => {
            console.log('WebSocket connection closed');
        };

        // Handle connection errors
        ws.onerror = (event) => {
            setError('WebSocket connection error');
            console.error('WebSocket error:', event);
        };

        // Cleanup: close connection on unmount or location change
        return () => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.close();
            }
        };
    }, [location]);

    return { drivers, error };
};