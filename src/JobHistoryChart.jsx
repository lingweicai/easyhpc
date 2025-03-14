import React, { useEffect, useState } from "react";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";
import cockpit from "cockpit";

const JobHistoryChart = () => {
        const [data, setData] = useState([]);

        useEffect(() => {
                cockpit.script("sacct -X --format=Start,End,JobID,State --parsable2").then(output => {
                        const jobs = output.split("\n").slice(1).map(line => {
                                const [start,, jobId, state] = line.split("|");
                                return { name: jobId, time: new Date(start).getTime(), state };
                        });
                        setData(jobs);
                });
        }, []);

        return (
                <ResponsiveContainer width="100%" height={300}>
                        <LineChart data={data}>
                                <XAxis dataKey="time" tickFormatter={t => new Date(t).toLocaleDateString()} />
                                <YAxis />
                                <Tooltip />
                                <Line type="monotone" dataKey="state" stroke="#82ca9d" />
                        </LineChart>
                </ResponsiveContainer>
        );
};

export default JobHistoryChart;
