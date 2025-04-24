"use client";

import * as React from "react";
import {
  ColumnDef,
  getCoreRowModel,
  useReactTable,
  flexRender,
} from "@tanstack/react-table";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useEffect, useState } from "react";
import { WithBot, Trade } from "@/utils/types";
import { getTradesForBot } from "@/utils/botData";

export const columns: ColumnDef<Trade>[] = [
  {
    accessorKey: "ticker",
    header: "Ticker",
    cell: ({ row }) => (
      <div className="font-medium">{row.getValue("ticker")}</div>
    ),
  },
  {
    accessorKey: "action",
    header: "Action",
    cell: ({ row }) => {
      const op = row.getValue("action") as Trade["action"];
      const color = op === "buy" ? "text-green-500" : "text-red-500";
      return <div className={`capitalize font-semibold ${color}`}>{op}</div>;
    },
  },
  {
    accessorKey: "numShares",
    header: "Shares",
    cell: ({ row }) => <div>{row.getValue("numShares")}</div>,
  },
  {
    accessorKey: "unitCost",
    header: "Unit Cost",
    cell: ({ row }) => {
      const unitCost = parseFloat(row.getValue("unitCost"));
      const formatted = new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
      }).format(unitCost);
      return <div className="text-left">{formatted}</div>;
    },
  },  
  {
    accessorKey: "time",
    header: "Time",
    cell: ({ row }) => {
      const dateStr = new Date(row.getValue("time")).toLocaleString();
      return <div className="text-sm text-muted-foreground">{dateStr}</div>;
    },
  },
];

export function TradeTable({ botId }: WithBot) {
  const [data, setData] = useState<Trade[]>([]);

  useEffect(() => {
    const get = async () => {
      if (!botId) return;
      const trades = await getTradesForBot(botId);
      setData(trades);
    };
    get();
  }, [botId, data]);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="w-full rounded-lg border bg-background text-foreground shadow-md">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead
                  key={header.id}
                  className="text-xs uppercase text-muted-foreground"
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows.length > 0 ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="text-center py-6">
                No data available
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
