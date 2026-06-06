import StatCards from "@/components/StatCards";
import ProjectTable from "@/components/ProjectTable";
import ProcessTimeline from "@/components/ProcessTimeline";
import StatusDonut from "@/components/StatusDonut";
import DetailPanel from "@/components/DetailPanel";

export default function DashboardPage() {
  return (
    <>
      <main className="scroll-thin min-w-0 flex-1 overflow-y-auto p-5">
        <div className="space-y-5">
          <StatCards />
          <ProjectTable />
          <div className="grid grid-cols-1 gap-5 xl:grid-cols-2">
            <ProcessTimeline />
            <StatusDonut />
          </div>
        </div>
      </main>

      <div className="hidden lg:block">
        <DetailPanel />
      </div>
    </>
  );
}
