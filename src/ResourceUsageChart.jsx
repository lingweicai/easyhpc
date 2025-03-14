import React, { useEffect, useState } from "react";
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from "recharts";
import cockpit from "cockpit";

const COLORS = ["#0088FE", "#00C49F"];

const ResourceUsageChart = () => {
        const [cpuUsage, setCpuUsage] = useState({ allocated: 0, idle: 0 });

        useEffect(() => {
                cockpit.script("sinfo -h -o '%C'").then(output => {
                        const [allocated, idle] = output.trim().split("/").map(Number);
                        setCpuUsage({ allocated, idle });
                });
        }, []);

        const data = [
                { name: "Allocated", value: cpuUsage.allocated },
                { name: "Idle", value: cpuUsage.idle },
        ];

        return (
                <ResponsiveContainer width="100%" height={300}>
                        <PieChart>
                                <Pie data={data} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={100}>
                                        {data.map((_, index) => (
                                                <Cell key={`cell-${index}`} fill={COLORS[index]} />
                                        ))}
                                </Pie>
                                <Tooltip />
                        </PieChart>
                </ResponsiveContainer>
        );
};

export default ResourceUsageChart;
