import { BadgeCheck, Heart, MessageCircle, Repeat2, Share2 } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

export type PostConfig = {
  text: string;
  name: string;
  handle: string;
  verified: boolean;
  avatarUrl: string;
  date: string;
  likeCount: string;
  cta: string;
  width: "tight" | "wide";
  mode: "classic" | "simple";
};

function getInitials(name: string) {
  const parts = name.trim().split(" ").filter(Boolean);
  if (parts.length === 0) return "X";
  const first = parts[0]?.[0] ?? "";
  const last = parts.length > 1 ? parts[parts.length - 1]?.[0] ?? "" : "";
  return `${first}${last}`.toUpperCase();
}

export function XPostPreview({ data }: { data: PostConfig }) {
  const handle = data.handle.trim();
  const formattedHandle = handle.startsWith("@") ? handle : `@${handle}`;
  const initials = getInitials(data.name);
  const showMetrics = data.mode === "classic";

  return (
    <div
      className={cn(
        "panel relative overflow-hidden px-6 py-6",
        data.width === "tight" ? "max-w-xl" : "max-w-2xl"
      )}
    >
      <div className="absolute left-0 top-0 h-1 w-full bg-gradient-to-r from-accent via-accent-2 to-transparent" />
      <div className="flex items-start gap-4">
        <div className="relative h-14 w-14 shrink-0">
          {data.avatarUrl ? (
            <img
              src={data.avatarUrl}
              alt="Avatar"
              className="h-full w-full rounded-2xl object-cover"
            />
          ) : (
            <div className="grid h-full w-full place-items-center rounded-2xl bg-black/5 text-sm font-semibold text-ink">
              {initials}
            </div>
          )}
          <span className="absolute -bottom-2 -right-2 rounded-full bg-white px-2 py-1 text-[10px] font-semibold text-accent shadow-soft">
            Live
          </span>
        </div>
        <div className="flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-base font-semibold text-ink">{data.name}</span>
            {data.verified ? (
              <Badge variant="default" className="gap-1">
                <BadgeCheck className="h-3.5 w-3.5" />
                Verified
              </Badge>
            ) : null}
          </div>
          <div className="text-sm text-muted">{formattedHandle}</div>
        </div>
      </div>

      <p className="mt-5 whitespace-pre-wrap text-base leading-relaxed text-ink">
        {data.text}
      </p>

      <div className="mt-4 flex flex-wrap items-center gap-3 text-xs text-muted">
        <span>{data.date}</span>
        {data.cta ? <span className="text-accent">{data.cta}</span> : null}
      </div>

      {showMetrics ? (
        <div className="mt-5 grid grid-cols-2 gap-3 text-xs text-muted sm:grid-cols-4">
          <div className="flex items-center gap-2">
            <MessageCircle className="h-4 w-4" />
            Reply
          </div>
          <div className="flex items-center gap-2">
            <Repeat2 className="h-4 w-4" />
            Repost
          </div>
          <div className="flex items-center gap-2">
            <Heart className="h-4 w-4" />
            {data.likeCount || "0"} likes
          </div>
          <div className="flex items-center gap-2">
            <Share2 className="h-4 w-4" />
            Share
          </div>
        </div>
      ) : null}
    </div>
  );
}
