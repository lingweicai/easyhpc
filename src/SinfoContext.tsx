// SinfoContext.tsx
import React, {
    createContext,
    useContext,
    useEffect,
    useState,
    useCallback,
} from 'react';
import cockpit from 'cockpit';

interface SinfoRawContextType {
    rawOutput: string;
    loading: boolean;
    error: string | null;
    refresh: () => void;
    refreshCounter: number;
  }

const SinfoContext = createContext<SinfoRawContextType>({
    rawOutput: '',
    loading: true,
    error: null,
    refresh: () => {},
    refreshCounter: 0, // 👈 NEW
});

export const useSinfoContext = () => useContext(SinfoContext);

const SinfoProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [rawOutput, setRawOutput] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [refreshCounter, setRefreshCounter] = useState(0);

    const fetchSinfo = useCallback(() => {
        setLoading(true);
        cockpit.spawn(['sinfo'])
                .then((output: string) => {
                    setRawOutput(output);
                    setError(null);
                })
                .catch((err: any) => {
                    console.error('sinfo fetch error:', err);
                    setError('Failed to fetch sinfo data');
                    setRawOutput('');
                })
                .finally(() => {
                    setLoading(false);
                    setRefreshCounter(prev => prev + 1);
                });
    }, []);

    useEffect(() => {
        fetchSinfo();
    }, []);

    return (
        <SinfoContext.Provider value={{ rawOutput, loading, error, refresh: fetchSinfo, refreshCounter }}>
            {children}
        </SinfoContext.Provider>
    );
};

export default SinfoProvider;
