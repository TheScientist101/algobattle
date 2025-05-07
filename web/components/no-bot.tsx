import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { NoBotPromptProps } from "@/utils/types";

/**
 * Displays a prompt encouraging users to create a bot if they don't have one.
 * - Opens a dialog modal to input and submit a bot name.
 * - On confirmation, calls `onCreate(name)` and closes the modal.
 *
 * Props:
 * - `onCreate`: function to call with the new bot name
 * - `loading`: whether the app is currently loading user/bot data
 */
export function NoBotPrompt({ onCreate }: NoBotPromptProps) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");

  /**
   * Handles creation logic:
   * - Trims whitespace
   * - Calls parent callback with name
   * - Clears input and closes modal
   */
  const handleCreate = () => {
    if (name.trim()) {
      onCreate(name);
      setName("");
      setOpen(false);
    }
  };

  return (
    <div className="flex items-center justify-center h-96">
      <div>
        {/* Dialog modal to create a bot */}
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild className="text-white text-lg">
            <Button variant="outline">Create a bot to get started!</Button>
          </DialogTrigger>

          <DialogContent className="sm:max-w-[425px] rounded-2xl bg-zinc-900 shadow-xl p-6">
            <DialogHeader>
              <DialogTitle className="text-white text-lg font-semibold">
                Create your bot
              </DialogTitle>
              <DialogDescription className="text-zinc-300">
                Enter the name of the bot
              </DialogDescription>
            </DialogHeader>

            {/* Input field for the name of the bot */}
            <div className="grid gap-4 py-4">
              <div className="flex flex-col gap-2">
                <label
                  htmlFor="name"
                  className="text-white text-sm font-medium"
                >
                  Bot Name
                </label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full bg-zinc-800 border border-zinc-700 text-white placeholder-zinc-500 focus:ring-primary focus:border-primary"
                  placeholder="Enter bot name"
                />
              </div>
            </div>

            {/* Save button */}
            <DialogFooter>
              <Button
                type="button"
                onClick={handleCreate}
                className="bg-primary text-white hover:bg-primary/90 transition-all"
              >
                Save
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
